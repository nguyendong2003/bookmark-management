package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	urlExpTime = 24 * time.Hour
)

type URLStorage interface {
	StoreURL(ctx context.Context, code, url string) error
}

type urlStorage struct {
	redisClient *redis.Client
}

func NewURLStorage(redisClient *redis.Client) URLStorage {
	return &urlStorage{
		redisClient: redisClient,
	}
}

func (s *urlStorage) StoreURL(ctx context.Context, code, url string) error {
	return s.redisClient.Set(ctx, code, url, urlExpTime).Err()
}
