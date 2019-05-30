package aws

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/stretchr/testify/mock"
)

// MockS3Svc mocks s3
type MockS3Svc struct {
	s3iface.S3API
	mock.Mock
}

// NewMockS3 mocks s3
func NewMockS3() *MockS3Svc {
	return &MockS3Svc{}
}

// ListBucketsWithContext lits
func (s *MockS3Svc) ListBucketsWithContext(ctx aws.Context, in *s3.ListBucketsInput, ro ...request.Option) (*s3.ListBucketsOutput, error) {
	args := s.Called(in)
	out := args.Get(0).(*s3.ListBucketsOutput)
	return out, args.Error(1)
}

// GetBucketLocationWithContext gets
func (s *MockS3Svc) GetBucketLocationWithContext(ctx aws.Context, in *s3.GetBucketLocationInput, ro ...request.Option) (*s3.GetBucketLocationOutput, error) {
	args := s.Called(in)
	out := args.Get(0).(*s3.GetBucketLocationOutput)
	return out, args.Error(1)
}

// GetBucketTaggingWithContext tags
func (s *MockS3Svc) GetBucketTaggingWithContext(ctx aws.Context, in *s3.GetBucketTaggingInput, ro ...request.Option) (*s3.GetBucketTaggingOutput, error) {
	args := s.Called(in)
	out := args.Get(0).(*s3.GetBucketTaggingOutput)
	return out, args.Error(1)
}

// GetBucketAclWithContext gets
func (s *MockS3Svc) GetBucketAclWithContext(ctx aws.Context, in *s3.GetBucketAclInput, ro ...request.Option) (*s3.GetBucketAclOutput, error) { // nolint: golint
	args := s.Called(in)
	out := args.Get(0).(*s3.GetBucketAclOutput)
	return out, args.Error(1)
}

// GetObjectWithContext gets
func (s *MockS3Svc) GetObjectWithContext(ctx aws.Context, in *s3.GetObjectInput, ro ...request.Option) (*s3.GetObjectOutput, error) {
	args := s.Called(in)
	out := args.Get(0).(*s3.GetObjectOutput)
	return out, args.Error(1)
}

// PutObjectWithContext put
func (s *MockS3Svc) PutObjectWithContext(ctx aws.Context, in *s3.PutObjectInput, ro ...request.Option) (*s3.PutObjectOutput, error) {
	args := s.Called(in)
	out := args.Get(0).(*s3.PutObjectOutput)
	return out, args.Error(1)
}

// MockS3Manager mocks s3manager
type MockS3Manager struct {
	s3manageriface.DownloaderAPI
	s3manageriface.UploaderAPI

	mock.Mock
}

// NewMockS3Manager returns a mock s3 manager
func NewMockS3Manager() *MockS3Manager {
	return &MockS3Manager{}
}

// DownloadWithContext downloads with context
func (s *MockS3Manager) DownloadWithContext(ctx aws.Context, writer io.WriterAt, in *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {
	args := s.Called(writer, in)
	out := args.Get(0).(int64)
	return out, args.Error(1)
}
