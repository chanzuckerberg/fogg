package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

//EC2 is an ec2 svc
type EC2 struct {
	Svc ec2iface.EC2API
}

// NewEC2 returns a new EC2 svc
func NewEC2(c client.ConfigProvider, config *aws.Config) *EC2 {
	return &EC2{Svc: ec2.New(c, config)}
}

// GetAllInstances will walk all instances and call func for each
func (e *EC2) GetAllInstances(ctx context.Context, f func(*ec2.Instance)) error {
	var err error
	input := &ec2.DescribeInstancesInput{}
	err = e.Svc.DescribeInstancesPagesWithContext(ctx, input, func(output *ec2.DescribeInstancesOutput, lastPage bool) bool {
		for _, reservation := range output.Reservations {
			if reservation == nil {
				continue
			}
			for _, instance := range reservation.Instances {
				f(instance)
			}
		}
		return true
	})

	return errors.Wrap(err, "error when getting all EC2 instances")
}

// GetAllVPCs will call f for each VPCs
func (e *EC2) GetAllVPCs(ctx context.Context, f func(*ec2.Vpc)) error {
	input := &ec2.DescribeVpcsInput{}
	out, err := e.Svc.DescribeVpcsWithContext(ctx, input)
	if err != nil {
		return err
	}
	for _, vpc := range out.Vpcs {
		f(vpc)
	}
	return nil
}
