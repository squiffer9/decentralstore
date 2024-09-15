package api_test

import (
	"bytes"
	"decentralstore/file-service/internal/api"
	"decentralstore/file-service/internal/domain"
	"decentralstore/file-service/internal/mocks"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileHandler_UploadFile(t *testing.T) {
	mockUseCase := new(mocks.MockFileUseCase)
	handler := api.NewFileHandler(mockUseCase)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	io.Copy(part, strings.NewReader("test content"))
	writer.Close()

	mockUseCase.On("UploadFile", mock.Anything, mock.Anything, "test.txt").Return(&domain.File{
		ID:   "123",
		Name: "test.txt",
		CID:  "QmTest123",
	}, nil)

	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.UploadFile(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, "123", response["id"])
	assert.Equal(t, "test.txt", response["name"])
	assert.Equal(t, "QmTest123", response["cid"])
	mockUseCase.AssertExpectations(t)
}

func TestFileHandler_DownloadFile(t *testing.T) {
	mockUseCase := new(mocks.MockFileUseCase)
	handler := api.NewFileHandler(mockUseCase)

	mockUseCase.On("DownloadFile", mock.Anything, "123", "test-keyword").Return(io.NopCloser(strings.NewReader("test content")), nil)

	req, _ := http.NewRequest("GET", "/download?id=123&keyword=test-keyword", nil)
	rr := httptest.NewRecorder()

	handler.DownloadFile(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test content", rr.Body.String())
	assert.Equal(t, "attachment; filename=123", rr.Header().Get("Content-Disposition"))
	assert.Equal(t, "application/octet-stream", rr.Header().Get("Content-Type"))
	mockUseCase.AssertExpectations(t)
}

func TestFileHandler_DeleteFile(t *testing.T) {
	mockUseCase := new(mocks.MockFileUseCase)
	handler := api.NewFileHandler(mockUseCase)

	mockUseCase.On("DeleteFile", mock.Anything, "123", "test-keyword").Return(nil)

	req, _ := http.NewRequest("DELETE", "/delete?id=123&keyword=test-keyword", nil)
	rr := httptest.NewRecorder()

	handler.DeleteFile(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "File deleted successfully", rr.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestFileHandler_DeleteFile_InvalidKeyword(t *testing.T) {
	mockUseCase := new(mocks.MockFileUseCase)
	handler := api.NewFileHandler(mockUseCase)

	mockUseCase.On("DeleteFile", mock.Anything, "123", "wrong-keyword").Return(&domain.ErrInvalidKeyword{Operation: "delete"})

	req, _ := http.NewRequest("DELETE", "/delete?id=123&keyword=wrong-keyword", nil)
	rr := httptest.NewRecorder()

	handler.DeleteFile(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid keyword for delete operation")
	mockUseCase.AssertExpectations(t)
}
