package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock implementation of redis.Client for testing
type MockRedisClient struct {
	mock.Mock
}

// Set implements the Set method of redis.Client
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

// Get implements the Get method of redis.Client
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

// Del implements the Del method of redis.Client
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func TestTokenCache_Set(t *testing.T) {
	ctx := context.Background()
	key := "test_key"
	value := "test_value"
	expiration := time.Hour

	// Test case 1: Successful Set
	mockRedis := new(MockRedisClient)
	mockRedis.On("Set", ctx, key, value, expiration).Return(redis.NewStatusCmd(ctx, "set", "OK")).Once()
	cache := NewTokenCache(mockRedis)
	err := cache.Set(ctx, key, value, expiration)
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)

	// Test case 2: Error during Set
	mockRedis = new(MockRedisClient)
	mockRedis.On("Set", ctx, key, value, expiration).Return(redis.NewStatusCmd(ctx, "set", errors.New("redis error"))).Once()
	cache = NewTokenCache(mockRedis)
	err = cache.Set(ctx, key, value, expiration)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis error")
	mockRedis.AssertExpectations(t)
}

func TestTokenCache_Get(t *testing.T) {
	ctx := context.Background()
	key := "test_key"
	expectedValue := "test_value"

	// Test case 1: Successful Get
	mockRedis := new(MockRedisClient)
	mockRedis.On("Get", ctx, key).Return(redis.NewStringCmd(ctx, "get", expectedValue)).Once()
	cache := NewTokenCache(mockRedis)
	value, err := cache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
	mockRedis.AssertExpectations(t)

	// Test case 2: Key not found
	mockRedis = new(MockRedisClient)
	mockRedis.On("Get", ctx, key).Return(redis.NewStringCmd(ctx, "get", "", redis.Nil)).Once()
	cache = NewTokenCache(mockRedis)
	value, err := cache.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Empty(t, value)
	mockRedis.AssertExpectations(t)

	// Test case 3: Error during Get
	mockRedis = new(MockRedisClient)
	mockRedis.On("Get", ctx, key).Return(redis.NewStringCmd(ctx, "get", "", errors.New("redis error"))).Once()
	cache = NewTokenCache(mockRedis)
	value, err := cache.Get(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis error")
	assert.Empty(t, value)
	mockRedis.AssertExpectations(t)
}

func TestTokenCache_Del(t *testing.T) {
	ctx := context.Background()
	key := "test_key"

	// Test case 1: Successful Del
	mockRedis := new(MockRedisClient)
	mockRedis.On("Del", ctx, []string{key}).Return(redis.NewIntCmd(ctx, "del", 1)).Once()
	cache := NewTokenCache(mockRedis)
	err := cache.Del(ctx, key)
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)

	// Test case 2: Error during Del
	mockRedis = new(MockRedisClient)
	mockRedis.On("Del", ctx, []string{key}).Return(redis.NewIntCmd(ctx, "del", errors.New("redis error"))).Once()
	cache = NewTokenCache(mockRedis)
	err = cache.Del(ctx, key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redis error")
	mockRedis.AssertExpectations(t)
}
