package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisBlacklistRepository implements the BlacklistRepository interface for Redis
type RedisBlacklistRepository struct {
	client *redis.Client
}

// NewRedisBlacklistRepository creates a new RedisBlacklistRepository
func NewRedisBlacklistRepository(client *redis.Client) *RedisBlacklistRepository {
	return &RedisBlacklistRepository{
		client: client,
	}
}

// Add adds a token to the blacklist with a given expiration time
func (r *RedisBlacklistRepository) Add(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	status := r.client.Set(ctx, key, true, expiration)
	if status.Err() != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", status.Err())
	}
	return nil
}

// Check checks if a token is present in the blacklist
func (r *RedisBlacklistRepository) Check(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token in blacklist: %w", err)
	}
	return val == 1, nil
}
