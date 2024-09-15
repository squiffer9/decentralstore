package infrastructure

import (
	"context"
	"math/big"
	"testing"

	"decentralstore/blockchain-service/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLatestBlockNumber(t *testing.T) {
	mockClient := new(mocks.MockEthClient)
	ethereumClient := &EthereumClient{client: mockClient}

	expectedBlockNumber := uint64(12345)
	mockClient.On("BlockNumber", mock.Anything).Return(expectedBlockNumber, nil)

	blockNumber, err := ethereumClient.GetLatestBlockNumber(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(int64(expectedBlockNumber)), blockNumber)
	mockClient.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	mockClient := new(mocks.MockEthClient)
	ethereumClient := &EthereumClient{client: mockClient}

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	expectedBalance := big.NewInt(1000000000000000000) // 1 ETH
	mockClient.On("BalanceAt", mock.Anything, address, mock.Anything).Return(expectedBalance, nil)

	balance, err := ethereumClient.GetBalance(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
	mockClient.AssertExpectations(t)
}

func TestSendTransaction(t *testing.T) {
	mockClient := new(mocks.MockEthClient)
	ethereumClient := &EthereumClient{client: mockClient}

	tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	mockClient.On("SendTransaction", mock.Anything, tx).Return(nil)

	err := ethereumClient.SendTransaction(context.Background(), tx)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWaitForTransaction(t *testing.T) {
	mockClient := new(mocks.MockEthClient)
	ethereumClient := &EthereumClient{client: mockClient}

	txHash := common.HexToHash("0x1234567890abcdef")
	expectedReceipt := &types.Receipt{TxHash: txHash, BlockNumber: big.NewInt(12345)}
	mockClient.On("TransactionReceipt", mock.Anything, txHash).Return(expectedReceipt, nil)

	receipt, err := ethereumClient.WaitForTransaction(context.Background(), txHash)

	assert.NoError(t, err)
	assert.Equal(t, expectedReceipt, receipt)
	mockClient.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	mockClient := new(mocks.MockEthClient)
	ethereumClient := &EthereumClient{client: mockClient}

	mockClient.On("Close").Return()

	ethereumClient.Close()

	mockClient.AssertExpectations(t)
}
