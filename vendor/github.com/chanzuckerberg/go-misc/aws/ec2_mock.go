package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/mock"
)

// MockEC2Svc is a mock of the ec2 service
type MockEC2Svc struct {
	ec2iface.EC2API
	mock.Mock
}

// NewMockEC2 returns a mock of ec2
func NewMockEC2() *MockEC2Svc {
	return &MockEC2Svc{}
}

// DescribeInstancesPagesWithContext is a mock
func (m *MockEC2Svc) DescribeInstancesPagesWithContext(ctx aws.Context, in *ec2.DescribeInstancesInput, fn func(*ec2.DescribeInstancesOutput, bool) bool, ro ...request.Option) error {
	args := m.Called(in)
	out := args.Get(0).(*ec2.DescribeInstancesOutput)
	err := args.Error(1)
	if err != nil {
		return err
	}
	fn(out, true)
	return nil
}
