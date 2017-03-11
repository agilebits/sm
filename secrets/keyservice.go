package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// EncryptionKey contians server key information
type EncryptionKey struct {
	KID    string `json:"kid"`
	Enc    string `json:"enc"`
	EncKey string `json:"encKey"`
	RawKey []byte `json:"-"`
}

// KeyService defines key methods
type KeyService interface {
	GenerateKey(kid string) (*EncryptionKey, error)
	DecryptKey(key *EncryptionKey) error
}

// Ciphertext contains encrypted message
type ciphertext struct {
	KID  string `json:"kid"`
	Enc  string `json:"enc"`
	Cty  string `json:"cty"`
	Iv   string `json:"iv"`
	Data string `json:"data"`
}

const (
	// A256GCM identifies the encryption algorithm
	A256GCM = "A256GCM"

	// B5JWKJSON identifies content type
	B5JWKJSON = "b5+jwk+json"
)

// Decrypt decrypts a given ciphertext byte array using the web crypto key
func (key *EncryptionKey) Decrypt(message []byte) ([]byte, error) {
	m := &ciphertext{}
	if err := json.Unmarshal(message, &m); err != nil {
		var errorMsg string
		if len(message) < 200 {
			errorMsg = fmt.Sprintf("invalid JSON %+q", string(message))
		} else {
			errorMsg = fmt.Sprintf("invalid JSON %+q...%+q", string(message[:120]), string(message[len(message)-75:]))
		}

		return nil, errors.Wrapf(err, errorMsg)
	}

	if m.KID != key.KID {
		return nil, fmt.Errorf("attempt to decrypt message with KID %v using different KID %v", m.KID, key.KID)
	}

	if m.Enc != A256GCM {
		return nil, fmt.Errorf("attempt to decrypt message with unknown enc: %+q", m.Enc)
	}

	if m.Cty != B5JWKJSON {
		return nil, fmt.Errorf("attempt to decrypt message with unknown cty: %+q", m.Cty)
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(m.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid ciphertext in the message")
	}

	block, err := aes.NewCipher(key.RawKey)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create NewCipher")
	}

	iv, err := base64.RawURLEncoding.DecodeString(m.Iv)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid iv in the message")
	}

	if len(iv) != 12 {
		return nil, fmt.Errorf("invalid iv length (%d) in the message, expected 12", len(iv))
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewGCM")
	}

	plaintext, err := aead.Open(nil, iv, ciphertext, nil)
	return plaintext, errors.Wrap(err, "failed to Open")
}

// Encrypt encrypts a given plaintext byte array
func (key *EncryptionKey) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.RawKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewCipher")
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create NewGCM")
	}

	iv := make([]byte, aead.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, errors.Wrap(err, "failed to get random iv")
	}

	data := aead.Seal(nil, iv, plaintext, nil)
	m := &ciphertext{
		KID:  key.KID,
		Enc:  A256GCM,
		Cty:  B5JWKJSON,
		Iv:   base64.RawURLEncoding.EncodeToString(iv),
		Data: base64.RawURLEncoding.EncodeToString(data),
	}

	result, err := json.Marshal(m)
	return result, errors.Wrap(err, "failed to Marshal")
}
