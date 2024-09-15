package domain

import (
	"math/big"
	"time"
)

// FileMetadata represents the metadata of a file stored on the blockchain
type FileMetadata struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Size            int64     `json:"size"`
	CID             string    `json:"cid"`
	UploadedAt      time.Time `json:"uploadedAt"`
	DownloadKeyword string    `json:"downloadKeyword"`
	DeleteKeyword   string    `json:"deleteKeyword"`
	Owner           string    `json:"owner"`
	BlockNumber     *big.Int  `json:"blockNumber"`
	TransactionHash string    `json:"transactionHash"`
}

// NewFileMetadata creates a new FileMetadata instance
func NewFileMetadata(id, name, cid, downloadKeyword, deleteKeyword, owner string, size int64) *FileMetadata {
	return &FileMetadata{
		ID:              id,
		Name:            name,
		Size:            size,
		CID:             cid,
		UploadedAt:      time.Now(),
		DownloadKeyword: downloadKeyword,
		DeleteKeyword:   deleteKeyword,
		Owner:           owner,
	}
}

// SetBlockchainInfo sets the blockchain-specific information for the metadata
func (fm *FileMetadata) SetBlockchainInfo(blockNumber *big.Int, transactionHash string) {
	fm.BlockNumber = blockNumber
	fm.TransactionHash = transactionHash
}
