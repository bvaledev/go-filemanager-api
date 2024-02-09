package domain

import (
	"context"
	"io"
)

type FileUploadAdapter interface {
	Upload(ctx context.Context, file io.Reader, filename string) error
}

type FileService interface {
	UploadFiles(files []FileInfo) []string
}
