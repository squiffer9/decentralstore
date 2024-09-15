package infrastructure

import (
	"context"
	"errors"
	"math/big"
	"time"

	"decentralstore/blockchain-service/internal/domain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type FileMetadataContract struct {
	address  common.Address
	backend  bind.ContractBackend
	contract boundContract
}

type boundContract interface {
	GetMetadata(opts *bind.CallOpts, fileID string) (struct {
		Name            string
		Size            uint64
		CID             string
		UploadedAt      uint64
		DownloadKeyword string
		DeleteKeyword   string
		Owner           common.Address
	}, error)
	StoreMetadata(opts *bind.TransactOpts, fileID string, name string, size uint64, cid string, downloadKeyword string, deleteKeyword string) (*types.Transaction, error)
	UpdateMetadata(opts *bind.TransactOpts, fileID string, isDeleted bool) (*types.Transaction, error)
}

func NewFileMetadataContract(address common.Address, backend bind.ContractBackend) (*FileMetadataContract, error) {
	contract, err := bindFileMetadataContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FileMetadataContract{
		address:  address,
		backend:  backend,
		contract: contract,
	}, nil
}

func (fmc *FileMetadataContract) StoreMetadata(ctx context.Context, metadata *domain.FileMetadata, opts *bind.TransactOpts) (*types.Transaction, error) {
	return fmc.contract.StoreMetadata(opts, metadata.ID, metadata.Name, uint64(metadata.Size), metadata.CID, metadata.DownloadKeyword, metadata.DeleteKeyword)
}

func (fmc *FileMetadataContract) GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error) {
	opts := &bind.CallOpts{Context: ctx}
	metadata, err := fmc.contract.GetMetadata(opts, fileID)
	if err != nil {
		return nil, err
	}

	return &domain.FileMetadata{
		ID:              fileID,
		Name:            metadata.Name,
		Size:            int64(metadata.Size),
		CID:             metadata.CID,
		UploadedAt:      time.Unix(int64(metadata.UploadedAt), 0),
		DownloadKeyword: metadata.DownloadKeyword,
		DeleteKeyword:   metadata.DeleteKeyword,
		Owner:           metadata.Owner.Hex(),
	}, nil
}

func (fmc *FileMetadataContract) UpdateMetadata(ctx context.Context, fileID string, isDeleted bool, opts *bind.TransactOpts) (*types.Transaction, error) {
	return fmc.contract.UpdateMetadata(opts, fileID, isDeleted)
}

func (fmc *FileMetadataContract) WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	backend, ok := fmc.backend.(bind.DeployBackend)
	if !ok {
		return nil, errors.New("backend does not implement bind.DeployBackend")
	}

	for {
		receipt, err := backend.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
			// 再試行
		}
	}
}

// bindFileMetadataContract binds a generic wrapper to an already deployed contract.
func bindFileMetadataContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (boundContract, error) {
	// 注意: この実装はモックまたはスタブです。
	// 実際の運用では、自動生成されたコントラクトバインディングを使用する必要があります。
	return &mockBoundContract{}, nil
}

// 以下は、boundContractのモック実装です。
// 実際の運用では、これらのメソッドは自動生成されたコードに置き換えられます。

type mockBoundContract struct{}

func (bc *mockBoundContract) GetMetadata(opts *bind.CallOpts, fileID string) (struct {
	Name            string
	Size            uint64
	CID             string
	UploadedAt      uint64
	DownloadKeyword string
	DeleteKeyword   string
	Owner           common.Address
}, error) {
	// モックの実装
	return struct {
		Name            string
		Size            uint64
		CID             string
		UploadedAt      uint64
		DownloadKeyword string
		DeleteKeyword   string
		Owner           common.Address
	}{
		Name:            "Mock File",
		Size:            1024,
		CID:             "QmMockCID",
		UploadedAt:      uint64(time.Now().Unix()),
		DownloadKeyword: "mockDownloadKey",
		DeleteKeyword:   "mockDeleteKey",
		Owner:           common.HexToAddress("0x1234567890123456789012345678901234567890"),
	}, nil
}

func (bc *mockBoundContract) StoreMetadata(opts *bind.TransactOpts, fileID string, name string, size uint64, cid string, downloadKeyword string, deleteKeyword string) (*types.Transaction, error) {
	// モックの実装
	return types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
}

func (bc *mockBoundContract) UpdateMetadata(opts *bind.TransactOpts, fileID string, isDeleted bool) (*types.Transaction, error) {
	// モックの実装
	return types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil), nil
}
