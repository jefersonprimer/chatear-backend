package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

// MockRedisClient is a mock implementation of redis.Client for testing
type MockRedisClient struct {
	SetFunc    func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	ExistsFunc func(ctx context.Context, keys ...string) *redis.IntCmd
}

// Set implements the Set method of redis.Client
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	return redis.NewStatusCmd(ctx, "set", key, value, expiration)
}

// Exists implements the Exists method of redis.Client
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if m.ExistsFunc != nil {
		return m.ExistsFunc(ctx, keys...)
	}
	return redis.NewIntCmd(ctx, "exists", keys)
}

func TestRedisBlacklistRepository_Add(t *testing.T) {
	ctx := context.Background()
	token := "test_token"
	expiration := time.Hour

	// Test case 1: Successful addition
	mockClient := &MockRedisClient{
		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			assert.Equal(t, fmt.Sprintf("blacklist:%s", token), key)
			assert.Equal(t, true, value)
			assert.Equal(t, expiration, expiration)
			return redis.NewStatusCmd(ctx, "set", "OK")
		},
	}
	repo := NewRedisBlacklistRepository(mockClient)
	err := repo.Add(ctx, token, expiration)
	assert.NoError(t, err)

	// Test case 2: Error during addition
	mockClient = &MockRedisClient{
		SetFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			return redis.NewStatusCmd(ctx, "set", errors.New("redis error"))
		},
	}
	repo = NewRedisBlacklistRepository(mockClient)
	err = repo.Add(ctx, token, expiration)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add token to blacklist")
}

func TestRedisBlacklistRepository_Check(t *testing.T) {
	ctx := context.Background()
	token := "test_token"

	// Test case 1: Token exists in blacklist
	mockClient := &MockRedisClient{
		ExistsFunc: func(ctx context.Context, keys ...string) *redis.IntCmd {
			assert.Contains(t, keys, fmt.Sprintf("blacklist:%s", token))
			return redis.NewIntCmd(ctx, "exists", 1)
		},
	}
	repo := NewRedisBlacklistRepository(mockClient)
	exists, err := repo.Check(ctx, token)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test case 2: Token does not exist in blacklist
	mockClient = &MockRedisClient{
		ExistsFunc: func(ctx context.Context, keys ...string) *redis.IntCmd {
			assert.Contains(t, keys, fmt.Sprintf("blacklist:%s", token))
			return redis.NewIntCmd(ctx, "exists", 0)
		},
	}
	repo = NewRedisBlacklistRepository(mockClient)
	exists, err = repo.Check(ctx, token)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test case 3: Error during check
	mockClient = &MockRedisClient{
		ExistsFunc: func(ctx context.Context, keys ...string) *redis.IntCmd {
			return redis.NewIntCmd(ctx, "exists", errors.New("redis error"))
		},
	}
	repo = NewRedisBlacklistRepository(mockClient)
	exists, err = repo.Check(ctx, token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check token in blacklist")
	assert.False(t, exists)
}