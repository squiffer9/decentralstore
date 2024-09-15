package usecase

import (
	"context"

	"decentralstore/blockchain-service/internal/domain"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockchainService interface {
	StoreMetadata(ctx context.Context, metadata *domain.FileMetadata) error
	GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error)
	UpdateMetadata(ctx context.Context, fileID string, isDeleted bool) error
}

type FileMetadataContractInterface interface {
	StoreMetadata(ctx context.Context, metadata *domain.FileMetadata, opts *bind.TransactOpts) (*types.Transaction, error)
	GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error)
	UpdateMetadata(ctx context.Context, fileID string, isDeleted bool, opts *bind.TransactOpts) (*types.Transaction, error)
	WaitForTransaction(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

type blockchainServiceImpl struct {
	contract FileMetadataContractInterface
}

func NewBlockchainService(contract FileMetadataContractInterface) BlockchainService {
	return &blockchainServiceImpl{contract: contract}
}

func (s *blockchainServiceImpl) StoreMetadata(ctx context.Context, metadata *domain.FileMetadata) error {
	tx, err := s.contract.StoreMetadata(ctx, metadata, nil)
	if err != nil {
		return err
	}
	_, err = s.contract.WaitForTransaction(ctx, tx.Hash())
	return err
}

func (s *blockchainServiceImpl) GetMetadata(ctx context.Context, fileID string) (*domain.FileMetadata, error) {
	return s.contract.GetMetadata(ctx, fileID)
}

func (s *blockchainServiceImpl) UpdateMetadata(ctx context.Context, fileID string, isDeleted bool) error {
	tx, err := s.contract.UpdateMetadata(ctx, fileID, isDeleted, nil)
	if err != nil {
		return err
	}
	_, err = s.contract.WaitForTransaction(ctx, tx.Hash())
	return err
}
