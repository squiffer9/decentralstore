package infrastructure

import (
	"context"
	"io"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/go-redis/redis/v8"
)

type IPFSShell interface {
	Add(r io.Reader, options ...shell.AddOpts) (string, error)
	Cat(path string) (io.ReadCloser, error)
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type StorageClient struct {
	IPFSShell   IPFSShell
	RedisClient RedisClient
}

func NewStorageClient(ipfsAPI, redisURL string) (*StorageClient, error) {
	ipfsShell := shell.NewShell(ipfsAPI)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	return &StorageClient{
		IPFSShell:   ipfsShell,
		RedisClient: redisClient,
	}, nil
}
