package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/stretchr/testify/mock"
)

// MockKMSSvc is a mock of the KMS service
type MockKMSSvc struct {
	kmsiface.KMSAPI
	mock.Mock
}

// NewMockKMS returns a new mock kms svc
func NewMockKMS() *MockKMSSvc {
	return &MockKMSSvc{}
}

// EncryptWithContext mocks Encrypt
func (k *MockKMSSvc) EncryptWithContext(ctx aws.Context, in *kms.EncryptInput, ro ...request.Option) (*kms.EncryptOutput, error) {
	args := k.Called(in)
	out := args.Get(0).(*kms.EncryptOutput)
	return out, args.Error(1)
}

// DecryptWithContext decrypts
func (k *MockKMSSvc) DecryptWithContext(ctx aws.Context, in *kms.DecryptInput, ro ...request.Option) (*kms.DecryptOutput, error) {
	args := k.Called(in)
	out := args.Get(0).(*kms.DecryptOutput)
	return out, args.Error(1)
}
