package mocks

import (
	"context"
	"io"
	"time"

	"decentralstore/file-service/internal/domain"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/go-redis/redis/v8"
)

// MockFileUseCase はFileUseCaseのモック実装です
type MockFileUseCase struct {
	UploadFileFn   func(ctx context.Context, file io.Reader, filename string) (*domain.File, error)
	DownloadFileFn func(ctx context.Context, fileID string, keyword string) (io.ReadCloser, error)
	DeleteFileFn   func(ctx context.Context, fileID string, keyword string) error
}

func (m *MockFileUseCase) UploadFile(ctx context.Context, file io.Reader, filename string) (*domain.File, error) {
	return m.UploadFileFn(ctx, file, filename)
}

func (m *MockFileUseCase) DownloadFile(ctx context.Context, fileID string, keyword string) (io.ReadCloser, error) {
	return m.DownloadFileFn(ctx, fileID, keyword)
}

func (m *MockFileUseCase) DeleteFile(ctx context.Context, fileID string, keyword string) error {
	return m.DeleteFileFn(ctx, fileID, keyword)
}

// MockIPFSShell はshell.Shellのモック実装です
type MockIPFSShell struct {
	AddFn func(r io.Reader, options ...shell.AddOpts) (string, error)
	CatFn func(path string) (io.ReadCloser, error)
}

func (m *MockIPFSShell) Add(r io.Reader, options ...shell.AddOpts) (string, error) {
	return m.AddFn(r, options...)
}

func (m *MockIPFSShell) Cat(path string) (io.ReadCloser, error) {
	return m.CatFn(path)
}

// MockRedisClient はredis.Clientのモック実装です
type MockRedisClient struct {
	SetFn func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	GetFn func(ctx context.Context, key string) *redis.StringCmd
	DelFn func(ctx context.Context, keys ...string) *redis.IntCmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.SetFn(ctx, key, value, expiration)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.GetFn(ctx, key)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return m.DelFn(ctx, keys...)
}
