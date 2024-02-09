package storage

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type S3ClientMock struct {
	mock.Mock
}

var _ s3Client = (*S3ClientMock)(nil)

func (m *S3ClientMock) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	call := m.Called(ctx, params, optFns)
	return call.Get(0).(*s3.PutObjectOutput), call.Error(1)
}

func TestUpload(t *testing.T) {
	fileTxt := strings.NewReader("Hello, gopher!")

	t.Run("should call put object once", func(t *testing.T) {
		s3ClientMock := new(S3ClientMock)
		s3ClientMock.On("PutObject", mock.Anything, mock.Anything, mock.Anything).Return(&s3.PutObjectOutput{}, nil)
		SUT := S3AdapterImpl{
			S3Client: s3ClientMock,
			S3Bucket: "gopher",
		}

		SUT.Upload(context.TODO(), fileTxt, "file.txt")

		s3ClientMock.AssertNumberOfCalls(t, "PutObject", 1)
	})

	t.Run("should upload file successfully", func(t *testing.T) {
		s3ClientMock := new(S3ClientMock)
		s3ClientMock.On("PutObject", mock.Anything, mock.Anything, mock.Anything).Return(&s3.PutObjectOutput{}, nil)
		SUT := S3AdapterImpl{
			S3Client: s3ClientMock,
			S3Bucket: "gopher",
		}

		err := SUT.Upload(context.TODO(), fileTxt, "file.txt")

		assert.Nil(t, err)
	})

	t.Run("should return error if upload fails", func(t *testing.T) {
		s3ClientMock := new(S3ClientMock)
		s3ClientMock.On("PutObject", mock.Anything, mock.Anything, mock.Anything).Return(&s3.PutObjectOutput{}, errors.New("any")).Once()
		SUT := S3AdapterImpl{
			S3Client: s3ClientMock,
			S3Bucket: "gopher",
		}

		err := SUT.Upload(context.TODO(), fileTxt, "file.txt")

		assert.NotNil(t, err)
	})
}
