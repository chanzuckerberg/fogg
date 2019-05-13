package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/stretchr/testify/mock"
)

// This is a mock for the IAM Svc - mock more functions here as needed

// MockIAMSvc is a mock of IAM service
type MockIAMSvc struct {
	iamiface.IAMAPI
	mock.Mock
}

// NewMockIAM returns a mock IAM SVC
func NewMockIAM() *MockIAMSvc {
	return &MockIAMSvc{}
}

// GetUserWithContext mocks getUserWithContext
func (i *MockIAMSvc) GetUserWithContext(ctx aws.Context, in *iam.GetUserInput, ro ...request.Option) (*iam.GetUserOutput, error) {
	args := i.Called(in)
	out := args.Get(0).(*iam.GetUserOutput)
	return out, args.Error(1)
}

// ListMFADevicesPagesWithContext lists
func (i *MockIAMSvc) ListMFADevicesPagesWithContext(ctx aws.Context, in *iam.ListMFADevicesInput, fn func(*iam.ListMFADevicesOutput, bool) bool, ro ...request.Option) error {
	args := i.Called(in)
	out := args.Get(0).(*iam.ListMFADevicesOutput)
	err := args.Error(1)
	if err != nil {
		return err
	}
	fn(out, true)
	return nil
}

// ListUsersPagesWithContext lists
func (i *MockIAMSvc) ListUsersPagesWithContext(ctx aws.Context, in *iam.ListUsersInput, fn func(*iam.ListUsersOutput, bool) bool, ro ...request.Option) error {
	args := i.Called(in)
	out := args.Get(0).(*iam.ListUsersOutput)
	err := args.Error(1)
	if err != nil {
		return err
	}
	fn(out, true)
	return nil
}

// GetLoginProfileWithContext gets
func (i *MockIAMSvc) GetLoginProfileWithContext(ctx aws.Context, in *iam.GetLoginProfileInput, ro ...request.Option) (*iam.GetLoginProfileOutput, error) {
	args := i.Called(in)
	out := args.Get(0).(*iam.GetLoginProfileOutput)
	return out, args.Error(1)
}
