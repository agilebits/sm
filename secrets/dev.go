package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// DevKeyService contains DevKeyService information
type DevKeyService struct {
	masterKey *EncryptionKey
}

const masterKeyID = "master-key"

// NewDevKeyService returns an empty DevKeyService object
func NewDevKeyService() *DevKeyService {
	createDataDir()

	result := &DevKeyService{}
	masterKey, err := result.GenerateKey(masterKeyID)
	if err != nil {
		log.Fatal("NewDevKeyService failed to generate master key")
	}

	result.masterKey = masterKey
	return result
}

func dataDir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("failed to obtain current user:", err)
	}

	return path.Join(user.HomeDir, ".sm")
}

func createDataDir() {
	if err := os.MkdirAll(dataDir(), 0700); err != nil {
		log.Fatal("failed to Mkdir:", err)
	}
}

func filepathForKeyID(kid string) string {
	filename := ""
	for _, ch := range strings.ToLower(kid) {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
			filename = filename + string(ch)
		}
	}

	return path.Join(dataDir(), filename)
}

func readKey(kid string) (*EncryptionKey, error) {
	buf, err := ioutil.ReadFile(filepathForKeyID(kid))
	if err != nil {
		return nil, errors.Wrapf(err, "readKey failed to ReadFile")
	}

	key := &EncryptionKey{}
	if err := json.Unmarshal(buf, &key); err != nil {
		return nil, errors.Wrapf(err, "readKey failed to Unmarshal")
	}

	return key, nil
}

func writeKey(key *EncryptionKey) error {
	buf, err := json.Marshal(key)
	if err != nil {
		return errors.Wrap(err, "writeKey failed to Marshal")
	}

	if err := ioutil.WriteFile(filepathForKeyID(key.KID), buf, 0700); err != nil {
		return errors.Wrap(err, "writeKey failed to WriteFile")
	}

	return nil
}

// GenerateKey generates a new server key
func (s *DevKeyService) GenerateKey(kid string) (*EncryptionKey, error) {
	result, err := readKey(kid)
	if err == nil {
		// key already exist,
		if err = s.DecryptKey(result); err != nil {
			return nil, errors.Wrap(err, "GenerateKey failed to DecryptKey")
		}
		return result, nil
	}

	rawKey := make([]byte, 32)
	if _, err := rand.Read(rawKey); err != nil {
		return nil, errors.Wrap(err, "GenerateKey failed to rand.Read")
	}

	var encKey string
	if kid == masterKeyID {
		// master key is stored in unencrypted
		encKey = base64.RawURLEncoding.EncodeToString(rawKey)
	} else {
		ciphertext, err := s.masterKey.Encrypt(rawKey)
		if err != nil {
			log.Fatal("GenerateKey failed to Encrypt with masterKey:", err)
		}
		encKey = base64.RawURLEncoding.EncodeToString(ciphertext)
	}

	result = &EncryptionKey{
		KID:    kid,
		Enc:    A256GCM,
		EncKey: encKey,
		RawKey: rawKey,
	}

	if err := writeKey(result); err != nil {
		log.Fatal("GenerateKey failed to writeKey:", err)
	}

	return result, nil
}

// DecryptKey decrypts the dev key
func (s *DevKeyService) DecryptKey(key *EncryptionKey) error {
	if key.RawKey != nil {
		return nil
	}

	encKey, err := base64.RawURLEncoding.DecodeString(key.EncKey)
	if err != nil {
		return errors.Wrap(err, "failed to decode base64url value")
	}

	if key.KID == masterKeyID {
		key.RawKey = encKey
	} else {
		plaintext, err := s.masterKey.Decrypt(encKey)
		if err != nil {
			return errors.Wrap(err, "failed to decrypt with master key")
		}
		key.RawKey = plaintext
	}

	return nil
}
