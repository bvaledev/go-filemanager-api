package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/bvaledev/go-filemanager-s3-chi/internal/domain"
)

type s3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

type S3AdapterImpl struct {
	S3Client s3Client
	S3Bucket string
}

func NewS3Adapter(S3Client s3Client, S3Bucket string) *S3AdapterImpl {
	return &S3AdapterImpl{S3Client, S3Bucket}
}

var _ domain.FileUploadAdapter = (*S3AdapterImpl)(nil)

func (s *S3AdapterImpl) Upload(ctx context.Context, file io.Reader, filename string) error {
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.S3Bucket),
		Key:    aws.String("go/" + filename),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	return err
}
