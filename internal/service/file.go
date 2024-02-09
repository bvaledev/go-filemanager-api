package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bvaledev/go-filemanager-s3-chi/internal/domain"
)

type FileServiceImpl struct {
	FileUploadAdapter domain.FileUploadAdapter
}

func NewFileService(FileUploadAdapter domain.FileUploadAdapter) *FileServiceImpl {
	return &FileServiceImpl{FileUploadAdapter}
}

func (s *FileServiceImpl) UploadFiles(files []domain.FileInfo) []string {
	maxUpload := 4
	uploadControl := make(chan struct{}, maxUpload)
	defer close(uploadControl)
	failedUpload := make(chan []string, len(files))
	defer close(failedUpload)

	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		uploadControl <- struct{}{}
		go s.upload(&wg, file, uploadControl, failedUpload)
	}

	wg.Wait()

	select {
	case fails := <-failedUpload:
		return fails
	default:
		return []string{}
	}
}

func (s *FileServiceImpl) upload(waitGroup *sync.WaitGroup, file domain.FileInfo, uploadControl <-chan struct{}, failedUpload chan []string) {
	defer func() {
		waitGroup.Done()
		<-uploadControl
	}()

	err := s.FileUploadAdapter.Upload(context.TODO(), file.File, file.Name)
	if err != nil {
		fmt.Printf("Upload failed %s.\n %s\n", file.Name, err.Error())
		select {
		case failed := <-failedUpload:
			failed = append(failed, file.Name)
			failedUpload <- failed
		default:
			failedUpload <- []string{file.Name}
		}
		return
	}
}
