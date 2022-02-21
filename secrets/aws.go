package secrets

import (
	"context"
	"encoding/base64"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/pkg/errors"
)

// AwsKeyService represents connection to Amazon Web Services KMS
type AwsKeyService struct {
	region      string
	masterKeyID string
	service     *kms.Client

	sync.RWMutex
}

// NewAwsKeyService creates a new AwsKeyService in given AWS region and with the given masterKey identifier.
func NewAwsKeyService(region string, masterKeyID string) *AwsKeyService {
	return &AwsKeyService{
		region:      region,
		masterKeyID: masterKeyID,
	}
}

func awsSession(region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
}

func (s *AwsKeyService) setup() error {
	s.Lock()
	defer s.Unlock()

	if s.service == nil {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(s.region))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}

		s.service = kms.NewFromConfig(cfg)
	}

	return nil
}

// GenerateKey generates a brand new ServerKey.
func (s *AwsKeyService) GenerateKey(kid string) (*EncryptionKey, error) {
	if err := s.setup(); err != nil {
		return nil, errors.Wrapf(err, "failed to setup")
	}

	input := &kms.GenerateDataKeyInput{
		EncryptionContext: map[string]string{"kid": kid},
		GrantTokens:       []string{"Encrypt", "Decrypt"},
		KeyId:             aws.String(s.masterKeyID),
		KeySpec:           "AES_256",
	}

	out, err := s.service.GenerateDataKey(context.TODO(), input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to GenerateDataKey")
	}

	result := &EncryptionKey{
		KID:    kid,
		Enc:    A256GCM,
		RawKey: out.Plaintext,
		EncKey: base64.RawURLEncoding.EncodeToString(out.CiphertextBlob),
	}

	return result, nil
}

// DecryptKey decrypts an existing ServerKey.
func (s *AwsKeyService) DecryptKey(key *EncryptionKey) error {
	if err := s.setup(); err != nil {
		return errors.Wrapf(err, "failed to setup")
	}

	ciphertextBlob, err := base64.RawURLEncoding.DecodeString(key.EncKey)
	if err != nil {
		return errors.Wrap(err, "failed to DecodeString")
	}

	input := &kms.DecryptInput{
		CiphertextBlob:    ciphertextBlob,
		EncryptionContext: map[string]string{"kid": key.KID},
		GrantTokens:       []string{"Encrypt", "Decrypt"},
	}

	out, err := s.service.Decrypt(context.TODO(), input)
	if err != nil {
		return errors.Wrapf(err, "failed to Decrypt")
	}

	key.RawKey = out.Plaintext
	return nil
}
