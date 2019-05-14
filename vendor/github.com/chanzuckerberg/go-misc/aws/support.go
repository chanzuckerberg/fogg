package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/support"
	"github.com/aws/aws-sdk-go/service/support/supportiface"
)

// Support is a support interface
type Support struct {
	Svc supportiface.SupportAPI
}

// NewSupport will return Support
func NewSupport(c client.ConfigProvider, config *aws.Config) *Support {
	return &Support{Svc: support.New(c, config)}
}
