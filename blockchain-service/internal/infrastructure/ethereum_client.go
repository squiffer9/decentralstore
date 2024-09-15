package infrastructure

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	client *ethclient.Client
}

func NewEthereumClient(url string) (*EthereumClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthereumClient{client: client}, nil
}

func (ec *EthereumClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return ec.client.CodeAt(ctx, contract, blockNumber)
}

func (ec *EthereumClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return ec.client.CallContract(ctx, call, blockNumber)
}

func (ec *EthereumClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	return ec.client.PendingCodeAt(ctx, account)
}

func (ec *EthereumClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return ec.client.PendingNonceAt(ctx, account)
}

func (ec *EthereumClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return ec.client.SuggestGasPrice(ctx)
}

func (ec *EthereumClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return ec.client.EstimateGas(ctx, call)
}

func (ec *EthereumClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return ec.client.SendTransaction(ctx, tx)
}

func (ec *EthereumClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	return ec.client.FilterLogs(ctx, query)
}

func (ec *EthereumClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	return ec.client.SubscribeFilterLogs(ctx, query, ch)
}

func (ec *EthereumClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return ec.client.TransactionReceipt(ctx, txHash)
}

// 以下のメソッドは既存のものですが、bind.ContractBackendインターフェースの完全な実装のために必要です

func (ec *EthereumClient) BlockNumber(ctx context.Context) (uint64, error) {
	return ec.client.BlockNumber(ctx)
}

func (ec *EthereumClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return ec.client.BalanceAt(ctx, account, blockNumber)
}

func (ec *EthereumClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return ec.client.HeaderByNumber(ctx, number)
}

func (ec *EthereumClient) Close() {
	ec.client.Close()
}

// SuggestGasTipCap is required for EIP-1559 transactions
func (ec *EthereumClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return ec.client.SuggestGasTipCap(ctx)
}

// Ensure EthereumClient implements bind.ContractBackend
var _ bind.ContractBackend = (*EthereumClient)(nil)
