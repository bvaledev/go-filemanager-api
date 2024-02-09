package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bvaledev/go-filemanager/internal/domain"
)

type FileServiceImpl struct {
	FileUploadAdapter domain.FileUploadAdapter
}

func (s *FileServiceImpl) UploadFiles(files []domain.FileInfo) []string {
	maxUpload := 4
	uploadControl := make(chan struct{}, maxUpload)
	defer close(uploadControl)
	failedUploads := make([]string, 0)
	failedUpload := make(chan string, maxUpload)

	go func() {
		defer close(failedUpload)
		for fileName := range failedUpload {
			failedUploads = append(failedUploads, fileName)
		}
	}()

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		uploadControl <- struct{}{}
		fmt.Println(file.Name)
		go s.upload(&wg, file, uploadControl, failedUpload)
	}
	wg.Wait()

	return failedUploads
}

func (s *FileServiceImpl) upload(waitGroup *sync.WaitGroup, file domain.FileInfo, uploadControl <-chan struct{}, failedUpload chan<- string) {
	defer func() {
		waitGroup.Done()
		<-uploadControl
	}()
	if err := s.FileUploadAdapter.Upload(context.TODO(), file.File, file.Name); err != nil {
		failedUpload <- file.Name
		return
	}
}
