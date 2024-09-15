package mocks

import (
	"context"
	"math/big"

	"decentralstore/blockchain-service/internal/domain"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

// MockEthClient is a mock of the EthClient interface
type MockEthClient struct {
	mock.Mock
}

func (m *MockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockEthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	args := m.Called(ctx, account, blockNumber)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	return args.Get(0).(*types.Receipt), args.Error(1)
}

func (m *MockEthClient) Close() {
	m.Called()
}

// MockContractBackend is a mock of the bind.ContractBackend interface
type MockContractBackend struct {
	mock.Mock
}

func (m *MockContractBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	args := m.Called(ctx, contract, blockNumber)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockContractBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	args := m.Called(ctx, call, blockNumber)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockContractBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*types.Header), args.Error(1)
}

func (m *MockContractBackend) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	args := m.Called(ctx, account)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockContractBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockContractBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockContractBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockContractBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	args := m.Called(ctx, call)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockContractBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockContractBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]types.Log), args.Error(1)
}

func (m *MockContractBackend) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	args := m.Called(ctx, query, ch)
	return args.Get(0).(ethereum.Subscription), args.Error(1)
}

func (m *MockContractBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	return args.Get(0).(*types.Receipt), args.Error(1)
}

// MockFileMetadataContract is a mock of the FileMetadataContractInterface
type MockFileMetadataContract struct {
	mock.Mock
}

func (m *MockFileMetadataContract) StoreMetadata(ctx context.Context, metadata *domain.FileMetadata, opts *bind.TransactOpts) (*types.Transaction, error) {
	args := m.Called(ctx, metadata, opts)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockFileMetadataContract) GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error) {
	args := m.Called(ctx, fileID)
	return args.Get(0).(*domain.FileMetadata), args.Error(1)
}

func (m *MockFileMetadataContract) UpdateMetadata(ctx context.Context, fileID string, isDeleted bool, opts *bind.TransactOpts) (*types.Transaction, error) {
	args := m.Called(ctx, fileID, isDeleted, opts)
	return args.Get(0).(*types.Transaction), args.Error(1)
}

func (m *MockFileMetadataContract) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := m.Called(ctx, txHash)
	return args.Get(0).(*types.Receipt), args.Error(1)
}
