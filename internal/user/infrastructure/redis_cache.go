package infrastructure

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// TokenRepository defines the interface for interacting with a token cache.
type TokenRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

// TokenCache is a Redis implementation of the TokenRepository.
type TokenCache struct {
	Client *redis.Client
}

// NewTokenCache creates a new TokenCache.
func NewTokenCache(redisURL string) (*TokenCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return &TokenCache{Client: client}, nil
}

// Set sets a key-value pair in the cache.
func (c *TokenCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Get gets a value from the cache by key.
func (c *TokenCache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// Del deletes a key from the cache.
func (c *TokenCache) Del(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}
