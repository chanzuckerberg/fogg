package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/pkg/errors"
)

// SecretsManager is a secretsmanager service
type SecretsManager struct {
	Svc secretsmanageriface.SecretsManagerAPI
}

// NewSecretsManager returns a new secrets manager
func NewSecretsManager(c client.ConfigProvider, config *aws.Config) *SecretsManager {
	return &SecretsManager{Svc: secretsmanager.New(c, config)}
}

// ReadStringLatestVersion reads the latest verison of a string secret
func (s *SecretsManager) ReadStringLatestVersion(ctx context.Context, secretID string) (*string, error) {
	input := &secretsmanager.GetSecretValueInput{}
	input.SetSecretId(secretID)
	output, err := s.Svc.GetSecretValueWithContext(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read secret with ID: %s", secretID)
	}
	if output == nil {
		return nil, errors.Wrapf(err, "Secret with ID %s is nil", secretID)
	}
	return output.SecretString, nil
}
