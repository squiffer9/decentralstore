package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"decentralstore/file-service/internal/infrastructure"
	"decentralstore/file-service/internal/mocks"
)

func TestSetupRoutes(t *testing.T) {
	mockIPFSShell := &mocks.MockIPFSShell{}
	mockRedisClient := &mocks.MockRedisClient{}
	
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFSShell,
		RedisClient: mockRedisClient,
	}

	router := SetupRoutes(storageClient)

	testServer := httptest.NewServer(router)
	defer testServer.Close()

	resp, err := http.Get(testServer.URL + "/upload")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed; got %v", resp.Status)
	}
}

func TestCreateFileHandler(t *testing.T) {
	mockIPFSShell := &mocks.MockIPFSShell{}
	mockRedisClient := &mocks.MockRedisClient{}
	
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFSShell,
		RedisClient: mockRedisClient,
	}

	handler := CreateFileHandler(storageClient)

	if handler == nil {
		t.Error("Expected non-nil FileHandler")
	}
}
