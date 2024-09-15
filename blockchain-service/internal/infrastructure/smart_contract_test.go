package infrastructure

import (
	"context"
	"math/big"
	"testing"
	"time"

	"decentralstore/blockchain-service/internal/domain"
	"decentralstore/blockchain-service/internal/mocks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func TestNewFileMetadataContract(t *testing.T) {
	mockBackend := new(mocks.MockContractBackend)
	address := common.HexToAddress("0x1234567890123456789012345678901234567890")

	contract, err := NewFileMetadataContract(address, mockBackend)

	assert.NoError(t, err)
	assert.NotNil(t, contract)
	assert.Equal(t, address, contract.address)
	assert.Equal(t, mockBackend, contract.backend)
}

func TestStoreMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	contract := &FileMetadataContract{
		address:  common.Address{},
		backend:  new(mocks.MockContractBackend),
		contract: mockContract,
	}

	ctx := context.Background()
	metadata := &domain.FileMetadata{
		ID:              "testID",
		Name:            "testName",
		Size:            1000,
		CID:             "testCID",
		DownloadKeyword: "testDownloadKeyword",
		DeleteKeyword:   "testDeleteKeyword",
	}
	opts := &bind.TransactOpts{}

	expectedTx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	mockContract.On("StoreMetadata", ctx, metadata, opts).Return(expectedTx, nil)

	tx, err := contract.StoreMetadata(ctx, metadata, opts)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
	mockContract.AssertExpectations(t)
}

func TestGetMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	contract := &FileMetadataContract{
		address:  common.Address{},
		backend:  new(mocks.MockContractBackend),
		contract: mockContract,
	}

	ctx := context.Background()
	fileID := "testID"
	expectedMetadata := &domain.FileMetadata{
		ID:              fileID,
		Name:            "testName",
		Size:            1000,
		CID:             "testCID",
		UploadedAt:      time.Now(),
		DownloadKeyword: "testDownloadKeyword",
		DeleteKeyword:   "testDeleteKeyword",
		Owner:           "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		BlockNumber:     big.NewInt(12345),
		TransactionHash: "0x1234567890abcdef",
	}

	mockContract.On("GetMetadata", ctx, fileID).Return(expectedMetadata, nil)

	metadata, err := contract.GetMetadata(ctx, fileID)

	assert.NoError(t, err)
	assert.Equal(t, expectedMetadata, metadata)
	mockContract.AssertExpectations(t)
}

func TestUpdateMetadata(t *testing.T) {
	mockContract := new(mocks.MockFileMetadataContract)
	contract := &FileMetadataContract{
		address:  common.Address{},
		backend:  new(mocks.MockContractBackend),
		contract: mockContract,
	}

	ctx := context.Background()
	fileID := "testID"
	isDeleted := true
	opts := &bind.TransactOpts{}
	expectedTx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)

	mockContract.On("UpdateMetadata", ctx, fileID, isDeleted, opts).Return(expectedTx, nil)

	tx, err := contract.UpdateMetadata(ctx, fileID, isDeleted, opts)

	assert.NoError(t, err)
	assert.Equal(t, expectedTx, tx)
	mockContract.AssertExpectations(t)
}

func TestWaitForTransactionInSmartContract(t *testing.T) {
	mockBackend := new(mocks.MockContractBackend)
	contract := &FileMetadataContract{
		address: common.Address{},
		backend: mockBackend,
	}

	ctx := context.Background()
	txHash := common.HexToHash("0x1234567890abcdef")
	expectedReceipt := &types.Receipt{TxHash: txHash, BlockNumber: big.NewInt(12345)}

	mockBackend.On("TransactionReceipt", ctx, txHash).Return(expectedReceipt, nil)

	receipt, err := contract.WaitForTransaction(ctx, txHash)

	assert.NoError(t, err)
	assert.Equal(t, expectedReceipt, receipt)
	mockBackend.AssertExpectations(t)
}
