package usecase_test

import (
	"context"
	"decentralstore/file-service/internal/domain"
	"decentralstore/file-service/internal/mocks"
	"decentralstore/file-service/internal/usecase"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileUseCaseImpl_UploadFile(t *testing.T) {
	mockStorage := new(mocks.MockStorageClient)
	fileUseCase := usecase.NewFileUseCase(mockStorage)

	ctx := context.Background()
	fileContent := "test file content"
	fileName := "test.txt"
	// cid := "QmTest123"

	mockStorage.On("StoreMetadata", mock.Anything, mock.AnythingOfType("*domain.File")).Return(nil)

	uploadedFile, err := fileUseCase.UploadFile(ctx, strings.NewReader(fileContent), fileName)

	assert.NoError(t, err)
	assert.NotNil(t, uploadedFile)
	assert.Equal(t, fileName, uploadedFile.Name)
	assert.NotEmpty(t, uploadedFile.CID)
	mockStorage.AssertExpectations(t)
}

func TestFileUseCaseImpl_DownloadFile(t *testing.T) {
	mockStorage := new(mocks.MockStorageClient)
	fileUseCase := usecase.NewFileUseCase(mockStorage)

	ctx := context.Background()
	fileID := "123"
	keyword := "test-keyword"
	cid := "QmTest123"

	mockFile := &domain.File{
		ID:              fileID,
		Name:            "test.txt",
		CID:             cid,
		DownloadKeyword: keyword,
	}

	mockStorage.On("GetMetadata", ctx, fileID).Return(mockFile, nil)

	reader, err := fileUseCase.DownloadFile(ctx, fileID, keyword)

	assert.NoError(t, err)
	assert.NotNil(t, reader)
	mockStorage.AssertExpectations(t)
}

func TestFileUseCaseImpl_DeleteFile(t *testing.T) {
	mockStorage := new(mocks.MockStorageClient)
	fileUseCase := usecase.NewFileUseCase(mockStorage)

	ctx := context.Background()
	fileID := "123"
	keyword := "test-keyword"

	mockFile := &domain.File{
		ID:            fileID,
		DeleteKeyword: keyword,
	}

	mockStorage.On("GetMetadata", ctx, fileID).Return(mockFile, nil)
	mockStorage.On("DeleteMetadata", ctx, fileID).Return(nil)

	err := fileUseCase.DeleteFile(ctx, fileID, keyword)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestFileUseCaseImpl_DeleteFile_InvalidKeyword(t *testing.T) {
	mockStorage := new(mocks.MockStorageClient)
	fileUseCase := usecase.NewFileUseCase(mockStorage)

	ctx := context.Background()
	fileID := "123"
	keyword := "wrong-keyword"

	mockFile := &domain.File{
		ID:            fileID,
		DeleteKeyword: "correct-keyword",
	}

	mockStorage.On("GetMetadata", ctx, fileID).Return(mockFile, nil)

	err := fileUseCase.DeleteFile(ctx, fileID, keyword)

	assert.Error(t, err)
	assert.IsType(t, &domain.ErrInvalidKeyword{}, err)
	mockStorage.AssertExpectations(t)
}
