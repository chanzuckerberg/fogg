package aws

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/pkg/errors"
)

// KMS is a kms client
type KMS struct {
	Svc kmsiface.KMSAPI
}

// NewKMS returns a KMS client
func NewKMS(s *session.Session, conf *aws.Config) *KMS {
	return &KMS{kms.New(s, conf)}
}

// EncryptBytes encrypts the plaintext using the keyID key and the given context
// result is base64 encoded string
func (k *KMS) EncryptBytes(ctx context.Context, keyID string, plaintext []byte, context map[string]*string) (string, error) {
	input := &kms.EncryptInput{}
	input.SetKeyId(keyID).SetPlaintext(plaintext).SetEncryptionContext(context)
	response, err := k.Svc.EncryptWithContext(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "KMS encryption failed")
	}
	if response == nil {
		return "", errors.New("Nil response from aws")
	}
	return base64.StdEncoding.EncodeToString(response.CiphertextBlob), nil
}

// Decrypt decrypts
func (k *KMS) Decrypt(ctx context.Context, ciphertext []byte, context map[string]*string) ([]byte, string, error) {
	input := &kms.DecryptInput{}
	input.SetCiphertextBlob(ciphertext).SetEncryptionContext(context)
	response, err := k.Svc.DecryptWithContext(ctx, input)
	if err != nil {
		return nil, "", errors.Wrap(err, "KMS decryption failed")
	}
	if response == nil {
		return nil, "", errors.New("Nil response from aws")
	}
	if response.KeyId == nil {
		return nil, "", errors.New("Nil KMS keyID returned from AWS")
	}
	return response.Plaintext, *response.KeyId, nil
}
