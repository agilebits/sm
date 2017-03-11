package secrets

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"

	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/service/kms"
)

// AwsKeyService represents connection to Amazon Web Services KMS
type AwsKeyService struct {
	lock sync.RWMutex

	region      string
	masterKeyID string

	creds   *credentials.Credentials
	service *kms.KMS
}

// NewAwsKeyService creates a new AwsKeyService in given AWS region and with the given masterKey identifier.
func NewAwsKeyService(region string, masterKeyID string) *AwsKeyService {
	return &AwsKeyService{
		region:      region,
		masterKeyID: masterKeyID,
	}
}

func awsSession(region string) *session.Session {
	return session.New(&aws.Config{
		Region: aws.String(region),
	})
}

func (s *AwsKeyService) setup() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.service == nil || s.creds == nil || s.creds.IsExpired() {
		s.creds = ec2rolecreds.NewCredentials(awsSession(s.region))
		sess := session.New(&aws.Config{
			Credentials: s.creds,
			Region:      &s.region,
		})

		s.service = kms.New(sess)
	}

	return nil
}

// GenerateKey generates a brand new ServerKey.
func (s *AwsKeyService) GenerateKey(kid string) (*EncryptionKey, error) {
	if err := s.setup(); err != nil {
		return nil, errors.Wrapf(err, "AwsKeyService.GenerateKey failed to setup")
	}

	input := &kms.GenerateDataKeyInput{
		EncryptionContext: map[string]*string{"kid": aws.String(kid)},
		GrantTokens:       []*string{aws.String("Encrypt"), aws.String("Decrypt")},
		KeyId:             aws.String(s.masterKeyID),
		KeySpec:           aws.String("AES_256"),
	}

	out, err := s.service.GenerateDataKey(input)
	if err != nil {
		return nil, errors.Wrapf(err, "AwsKeyService.GenerateKey failed to GenerateDataKey")
	}

	result := &EncryptionKey{
		KID:    *out.KeyId,
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
		EncryptionContext: map[string]*string{"kid": aws.String(key.KID)},
		GrantTokens:       []*string{aws.String("Encrypt"), aws.String("Decrypt")},
	}

	out, err := s.service.Decrypt(input)
	if err != nil {
		return errors.Wrapf(err, "failed to Decrypt")
	}

	key.RawKey = out.Plaintext
	return nil
}
