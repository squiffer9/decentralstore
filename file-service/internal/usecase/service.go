package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"decentralstore/file-service/internal/domain"
	"decentralstore/file-service/internal/infrastructure"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, file io.Reader, filename string) (*domain.File, error)
	DownloadFile(ctx context.Context, fileID string, keyword string) (io.ReadCloser, error)
	DeleteFile(ctx context.Context, fileID string, keyword string) error
}

type FileUseCaseImpl struct {
	StorageClient *infrastructure.StorageClient
}

func NewFileUseCase(storageClient *infrastructure.StorageClient) FileUseCase {
	return &FileUseCaseImpl{StorageClient: storageClient}
}

func (s *FileUseCaseImpl) UploadFile(ctx context.Context, file io.Reader, filename string) (*domain.File, error) {
	// IPFSにファイルをアップロード
	cid, err := s.StorageClient.IPFSShell.Add(file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to IPFS: %w", err)
	}

	// メタデータを作成
	id := generateUniqueID()
	downloadKeyword := generateKeyword()
	deleteKeyword := generateKeyword()
	uploadedFile := &domain.File{
		ID:              id,
		Name:            filename,
		CID:             cid,
		UploadedAt:      time.Now(),
		DownloadKeyword: downloadKeyword,
		DeleteKeyword:   deleteKeyword,
	}

	// Redisにメタデータを保存
	err = s.storeMetadata(ctx, uploadedFile)
	if err != nil {
		return nil, fmt.Errorf("failed to store metadata: %w", err)
	}

	return uploadedFile, nil
}

func (s *FileUseCaseImpl) DownloadFile(ctx context.Context, fileID string, keyword string) (io.ReadCloser, error) {
	// Redisからメタデータを取得
	metadata, err := s.getMetadata(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	// キーワードを検証
	if !validateKeyword(keyword, metadata.DownloadKeyword) {
		return nil, fmt.Errorf("invalid download keyword")
	}

	// IPFSからファイルを取得
	reader, err := s.StorageClient.IPFSShell.Cat(metadata.CID)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from IPFS: %w", err)
	}

	return reader, nil
}

func (s *FileUseCaseImpl) DeleteFile(ctx context.Context, fileID string, keyword string) error {
	// Redisからメタデータを取得
	metadata, err := s.getMetadata(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	// キーワードを検証
	if !validateKeyword(keyword, metadata.DeleteKeyword) {
		return fmt.Errorf("invalid delete keyword")
	}

	// Redisからメタデータを削除
	err = s.deleteMetadata(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}

func (s *FileUseCaseImpl) storeMetadata(ctx context.Context, file *domain.File) error {
	jsonData, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("failed to marshal file metadata: %w", err)
	}

	err = s.StorageClient.RedisClient.Set(ctx, "file:"+file.ID, string(jsonData), 0).Err()
	if err != nil {
		return fmt.Errorf("failed to store metadata in Redis: %w", err)
	}

	return nil
}

func (s *FileUseCaseImpl) getMetadata(ctx context.Context, fileID string) (*domain.File, error) {
	jsonData, err := s.StorageClient.RedisClient.Get(ctx, "file:"+fileID).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata from Redis: %w", err)
	}

	var file domain.File
	err = json.Unmarshal([]byte(jsonData), &file)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file metadata: %w", err)
	}

	return &file, nil
}

func (s *FileUseCaseImpl) deleteMetadata(ctx context.Context, fileID string) error {
	err := s.StorageClient.RedisClient.Del(ctx, "file:"+fileID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete metadata from Redis: %w", err)
	}

	return nil
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateKeyword() string {
	// 実際の実装では、より安全なランダム文字列生成メソッドを使用してください
	return fmt.Sprintf("key-%d", time.Now().UnixNano())
}

func validateKeyword(provided, stored string) bool {
	return provided == stored
}
