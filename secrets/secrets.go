package secrets

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const (
	envDev = "dev"
	envAWS = "aws"
)

// EncryptEnvelope will generate a new key and encrypt the message. It returns the Envelope that contains everything that is needed to decrypt the message (if the access to the KeyService is granted).
func EncryptEnvelope(env, region, masterKeyID string, message []byte) (*Envelope, error) {
	keyService, err := getKeyService(env, region, masterKeyID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain key service for env:%+q, region:%+q, masterKeyID:%+q", env, region, masterKeyID)
	}

	kid := "sm-" + time.Now().Format(time.RFC3339)
	encryptionKey, err := keyService.GenerateKey(kid)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate encryption key")
	}

	ciphertext, err := encryptionKey.Encrypt(message)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encrypt")
	}

	envelope := &Envelope{
		Env:         env,
		Region:      region,
		MasterKeyID: masterKeyID,
		Key:         *encryptionKey,
		Data:        base64.RawURLEncoding.EncodeToString(ciphertext),
	}

	return envelope, nil
}

// DecryptEnvelope will access the key service and decrypt the envelope.
func DecryptEnvelope(envelope *Envelope) ([]byte, error) {
	keyService, err := getKeyService(envelope.Env, envelope.Region, envelope.MasterKeyID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to obtain key service for env:%+q, region:%+q, masterKeyID:%+q", envelope.Env, envelope.Region, envelope.MasterKeyID)
	}

	encryptionKey := &envelope.Key
	if err := keyService.DecryptKey(encryptionKey); err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt key")
	}

	data, err := base64.RawURLEncoding.DecodeString(envelope.Data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode data")
	}

	result, err := encryptionKey.Decrypt(data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decrypt data")
	}

	return result, nil
}

func getKeyService(env, region, masterKeyID string) (KeyService, error) {
	var keyService KeyService
	switch env {
	case envDev:
		keyService = NewDevKeyService()
	case envAWS:
		keyService = NewAwsKeyService(region, masterKeyID)
	default:
		return nil, fmt.Errorf("unsupported env: %+q", env)
	}

	return keyService, nil
}
