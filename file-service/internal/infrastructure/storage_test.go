package infrastructure_test

import (
	"context"
	"decentralstore/file-service/internal/domain"
	"decentralstore/file-service/internal/infrastructure"
	"decentralstore/file-service/internal/mocks"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStorageClient_StoreMetadata(t *testing.T) {
	mockRedis := new(mocks.MockRedisClient)
	mockIPFS := new(mocks.MockIPFSShell)
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFS,
		RedisClient: mockRedis,
	}

	ctx := context.Background()
	file := &domain.File{
		ID:   "123",
		Name: "test.txt",
		CID:  "QmTest123",
	}

	mockRedis.On("Set", ctx, mock.Anything, mock.Anything, time.Duration(0)).Return(redis.NewStatusResult("OK", nil))

	err := storageClient.StoreMetadata(ctx, file)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestStorageClient_GetMetadata(t *testing.T) {
	mockRedis := new(mocks.MockRedisClient)
	mockIPFS := new(mocks.MockIPFSShell)
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFS,
		RedisClient: mockRedis,
	}

	ctx := context.Background()
	fileID := "123"
	fileJSON := `{"id":"123","name":"test.txt","cid":"QmTest123"}`

	mockRedis.On("Get", ctx, "file:"+fileID).Return(redis.NewStringResult(fileJSON, nil))

	file, err := storageClient.GetMetadata(ctx, fileID)

	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, fileID, file.ID)
	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, "QmTest123", file.CID)
	mockRedis.AssertExpectations(t)
}

func TestStorageClient_DeleteMetadata(t *testing.T) {
	mockRedis := new(mocks.MockRedisClient)
	mockIPFS := new(mocks.MockIPFSShell)
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFS,
		RedisClient: mockRedis,
	}

	ctx := context.Background()
	fileID := "123"

	mockRedis.On("Del", ctx, "file:"+fileID).Return(redis.NewIntResult(1, nil))

	err := storageClient.DeleteMetadata(ctx, fileID)

	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestStorageClient_GetMetadata_NotFound(t *testing.T) {
	mockRedis := new(mocks.MockRedisClient)
	mockIPFS := new(mocks.MockIPFSShell)
	storageClient := &infrastructure.StorageClient{
		IPFSShell:   mockIPFS,
		RedisClient: mockRedis,
	}

	ctx := context.Background()
	fileID := "123"

	mockRedis.On("Get", ctx, "file:"+fileID).Return(redis.NewStringResult("", redis.Nil))

	file, err := storageClient.GetMetadata(ctx, fileID)

	assert.Error(t, err)
	assert.Nil(t, file)
	assert.IsType(t, &domain.ErrNotFound{}, err)
	mockRedis.AssertExpectations(t)
}
