package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"context"

	"decentralstore/blockchain-service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlockchainService is a mock of BlockchainService interface
type MockBlockchainService struct {
	mock.Mock
}

func (m *MockBlockchainService) StoreMetadata(ctx context.Context, metadata *domain.FileMetadata) error {
	args := m.Called(ctx, metadata)
	return args.Error(0)
}

func (m *MockBlockchainService) GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(*domain.FileMetadata), args.Error(1)
}

func (m *MockBlockchainService) UpdateMetadata(ctx context.Context, fileID string, isDeleted bool) error {
	args := m.Called(ctx, fileID, isDeleted)
	return args.Error(0)
}

func TestStoreMetadata(t *testing.T) {
	mockService := new(MockBlockchainService)
	handler := NewBlockchainHandler(mockService)

	metadata := domain.FileMetadata{
		ID:   "testID",
		Name: "testFile",
		CID:  "QmTest",
	}

	mockService.On("StoreMetadata", mock.Anything, &metadata).Return(nil)

	body, _ := json.Marshal(metadata)
	req, _ := http.NewRequest("POST", "/store", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.StoreMetadata(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)
}

func TestGetMetadata(t *testing.T) {
	mockService := new(MockBlockchainService)
	handler := NewBlockchainHandler(mockService)

	fileID := "testID"
	expectedMetadata := &domain.FileMetadata{
		ID:   fileID,
		Name: "testFile",
		CID:  "QmTest",
	}

	mockService.On("GetMetadata", mock.Anything, fileID).Return(expectedMetadata, nil)

	req, _ := http.NewRequest("GET", "/metadata?fileID="+fileID, nil)
	rr := httptest.NewRecorder()

	handler.GetMetadata(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseMetadata domain.FileMetadata
	json.Unmarshal(rr.Body.Bytes(), &responseMetadata)
	assert.Equal(t, expectedMetadata, &responseMetadata)
	mockService.AssertExpectations(t)
}

func TestUpdateMetadata(t *testing.T) {
	mockService := new(MockBlockchainService)
	handler := NewBlockchainHandler(mockService)

	fileID := "testID"
	isDeleted := true

	mockService.On("UpdateMetadata", mock.Anything, fileID, isDeleted).Return(nil)

	updateRequest := map[string]bool{"isDeleted": isDeleted}
	body, _ := json.Marshal(updateRequest)
	req, _ := http.NewRequest("PUT", "/update?fileID="+fileID, bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.UpdateMetadata(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockService.AssertExpectations(t)
}

func TestMethodNotAllowed(t *testing.T) {
	mockService := new(MockBlockchainService)
	handler := NewBlockchainHandler(mockService)

	testCases := []struct {
		name        string
		method      string
		handlerFunc func(http.ResponseWriter, *http.Request)
	}{
		{"StoreMetadata", "GET", handler.StoreMetadata},
		{"GetMetadata", "POST", handler.GetMetadata},
		{"UpdateMetadata", "GET", handler.UpdateMetadata},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/", nil)
			rr := httptest.NewRecorder()

			tc.handlerFunc(rr, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
		})
	}
}
