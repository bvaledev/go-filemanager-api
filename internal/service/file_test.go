package service

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/bvaledev/go-filemanager-s3-chi/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type S3AdapterMock struct {
	mock.Mock
}

var _ domain.FileUploadAdapter = (*S3AdapterMock)(nil)

func (m *S3AdapterMock) Upload(ctx context.Context, file io.Reader, filename string) error {
	call := m.Called(ctx, file, filename)
	return call.Error(0)
}

func TestUploadFiles(t *testing.T) {
	fileTxt := strings.NewReader("Hello, gopher!")
	files := []domain.FileInfo{
		{
			Name: "file_1.txt",
			File: fileTxt,
		},
		{
			Name: "file_2.txt",
			File: fileTxt,
		},
		{
			Name: "file_3.txt",
			File: fileTxt,
		},
	}

	t.Run("should upload files successfully", func(t *testing.T) {
		s3AdapterMock := new(S3AdapterMock)
		s3AdapterMock.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		SUT := NewFileService(s3AdapterMock)

		failed := SUT.UploadFiles(files)

		assert.Equal(t, len(failed), 0)
	})

	t.Run("should return all failed files upload", func(t *testing.T) {
		s3AdapterMock := new(S3AdapterMock)
		s3AdapterMock.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("fail"))
		SUT := NewFileService(s3AdapterMock)

		failed := SUT.UploadFiles(files)

		assert.Equal(t, 3, len(failed))
	})
}

func Test_upload(t *testing.T) {
	fileTxt := strings.NewReader("Hello, gopher!")
	fileInfo := domain.FileInfo{
		Name: "file.txt",
		File: fileTxt,
	}
	s3AdapterMock := new(S3AdapterMock)

	t.Run("should call s3 upload with correct values", func(t *testing.T) {
		s3AdapterMock.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		SUT := NewFileService(s3AdapterMock)

		var awaitGroup sync.WaitGroup
		uploadControl := make(chan struct{}, 50)
		errorFileUpload := make(chan []string, 1)
		defer close(uploadControl)
		defer close(errorFileUpload)

		uploadControl <- struct{}{}
		awaitGroup.Add(1)
		go SUT.upload(&awaitGroup, fileInfo, uploadControl, errorFileUpload)
		awaitGroup.Wait()

		s3AdapterMock.AssertNumberOfCalls(t, "Upload", 1)
		s3AdapterMock.AssertCalled(t, "Upload", context.TODO(), fileInfo.File, "file.txt")
	})

	t.Run("should write error when s3 upload fails", func(t *testing.T) {
		s3AdapterMock.On("Upload", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("any")).Once()
		SUT := NewFileService(s3AdapterMock)

		var awaitGroup sync.WaitGroup
		uploadControl := make(chan struct{}, 50)
		errorFileUpload := make(chan []string, 1)
		defer close(uploadControl)
		defer close(errorFileUpload)

		uploadControl <- struct{}{}
		awaitGroup.Add(1)
		go SUT.upload(&awaitGroup, fileInfo, uploadControl, errorFileUpload)
		awaitGroup.Wait()
		failed := <-errorFileUpload
		assert.Equal(t, failed[0], "file.txt")
	})
}
