package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/pkg/errors"
)

// S3 is an s3 client
type S3 struct {
	Svc        s3iface.S3API
	Downloader s3manageriface.DownloaderAPI
}

// NewS3 returns an s3 client
func NewS3(c client.ConfigProvider, config *aws.Config) *S3 {
	client := s3.New(c, config)
	downloder := s3manager.NewDownloaderWithClient(client)
	return &S3{Svc: client, Downloader: downloder}
}

// ListBuckets lists buckets
func (s *S3) ListBuckets(ctx context.Context) (*s3.ListBucketsOutput, error) {
	input := &s3.ListBucketsInput{}
	out, err := s.Svc.ListBucketsWithContext(ctx, input)
	return out, errors.Wrap(err, "Error listing s3 buckets")
}

// GetBucketLocation gets the bucket's location (region)
func (s *S3) GetBucketLocation(ctx context.Context, bucketName string) (string, error) {
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	}

	out, err := s.Svc.GetBucketLocationWithContext(ctx, input)
	if err != nil {
		return "", errors.Wrapf(err, "Error getting bucket %s location", bucketName)
	}
	if out == nil {
		return "", errors.New("Nil output from aws")
	}

	// If nil then us-east-1
	if out.LocationConstraint == nil {
		return "us-east-1", nil
	}
	return *out.LocationConstraint, nil
}

// GetBucketTagging returns the bucket's tags
func (s *S3) GetBucketTagging(ctx context.Context, bucketName string) (*s3.GetBucketTaggingOutput, error) {
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}
	out, err := s.Svc.GetBucketTaggingWithContext(ctx, input)
	return out, errors.Wrapf(err, "Error getting bucket tags for %s", bucketName)
}

// GetBucketACL gets the bucket's ACL
func (s *S3) GetBucketACL(ctx context.Context, bucketName string) (*s3.GetBucketAclOutput, error) {
	input := &s3.GetBucketAclInput{}
	input.SetBucket(bucketName)

	out, err := s.Svc.GetBucketAclWithContext(ctx, input)
	return out, errors.Wrapf(err, "Error getting bucket %s ACL", bucketName)
}
