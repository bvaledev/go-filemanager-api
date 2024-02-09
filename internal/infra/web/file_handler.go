package web

import (
	"encoding/json"
	"net/http"

	"github.com/bvaledev/go-filemanager-s3-chi/internal/domain"
)

type FileHandler struct {
	FileService domain.FileService
}

func NewFileHandler(FileService domain.FileService) *FileHandler {
	return &FileHandler{FileService}
}

type output struct {
	Failed []string `json:"failed_uploads"`
}

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	// maximum upload of 40 MB.
	if err := r.ParseMultipartForm(40 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filesToUpload := make([]domain.FileInfo, 0)
	for _, fileHeader := range r.MultipartForm.File["files"] {
		file, _ := fileHeader.Open()
		defer file.Close()
		fileInfo := domain.FileInfo{
			Name: fileHeader.Filename,
			File: file,
		}
		filesToUpload = append(filesToUpload, fileInfo)
	}

	failed := h.FileService.UploadFiles(filesToUpload)

	if err := json.NewEncoder(w).Encode(output{Failed: failed}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
