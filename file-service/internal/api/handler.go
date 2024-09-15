package api

import (
	"io"
	"encoding/json"
	"net/http"

	"decentralstore/file-service/internal/usecase"
)

type FileHandler struct {
	fileUseCase usecase.FileUseCase
}

func NewFileHandler(fileUseCase usecase.FileUseCase) *FileHandler {
	return &FileHandler{fileUseCase: fileUseCase}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadedFile, err := h.fileUseCase.UploadFile(r.Context(), file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadedFile)
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := r.URL.Query().Get("id")
	keyword := r.URL.Query().Get("keyword")

	if fileID == "" || keyword == "" {
		http.Error(w, "Missing file ID or keyword", http.StatusBadRequest)
		return
	}

	reader, err := h.fileUseCase.DownloadFile(r.Context(), fileID, keyword)
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileID)
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, reader)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileID := r.URL.Query().Get("id")
	keyword := r.URL.Query().Get("keyword")

	if fileID == "" || keyword == "" {
		http.Error(w, "Missing file ID or keyword", http.StatusBadRequest)
		return
	}

	err := h.fileUseCase.DeleteFile(r.Context(), fileID, keyword)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}
