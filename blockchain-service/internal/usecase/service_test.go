package usecase

import (
	"context"
	"math/big"
	"testing"
	"time"

	"decentralstore/blockchain-service/internal/domain"
	"decentralstore/blockchain-service/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStoreMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	service := NewBlockchainService(mockContract)

	ctx := context.Background()
	metadata := &domain.FileMetadata{
		ID:              "testID",
		Name:            "testFile",
		Size:            1000,
		CID:             "QmTest",
		UploadedAt:      time.Now(),
		DownloadKeyword: "downloadKey",
		DeleteKeyword:   "deleteKey",
		Owner:           "0x1234567890123456789012345678901234567890",
	}

	mockTx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	mockReceipt := &types.Receipt{Status: 1}

	mockContract.On("StoreMetadata", ctx, metadata, mock.Anything).Return(mockTx, nil)
	mockContract.On("WaitForTransaction", ctx, mockTx.Hash()).Return(mockReceipt, nil)

	err := service.StoreMetadata(ctx, metadata)

	assert.NoError(t, err)
	mockContract.AssertExpectations(t)
}

func TestGetMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	service := NewBlockchainService(mockContract)

	ctx := context.Background()
	fileID := "testID"
	expectedMetadata := &domain.FileMetadata{
		ID:   fileID,
		Name: "testFile",
		Size: 1000,
		CID:  "QmTest",
	}

	mockContract.On("GetMetadata", ctx, fileID).Return(expectedMetadata, nil)

	metadata, err := service.GetMetadata(ctx, fileID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetadata, metadata)
	mockContract.AssertExpectations(t)
}

func TestUpdateMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	service := NewBlockchainService(mockContract)

	ctx := context.Background()
	fileID := "testID"
	isDeleted := true

	mockTx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	mockReceipt := &types.Receipt{Status: 1}

	mockContract.On("UpdateMetadata", ctx, fileID, isDeleted, mock.Anything).Return(mockTx, nil)
	mockContract.On("WaitForTransaction", ctx, mockTx.Hash()).Return(mockReceipt, nil)

	err := service.UpdateMetadata(ctx, fileID, isDeleted)

	assert.NoError(t, err)
	mockContract.AssertExpectations(t)
}
