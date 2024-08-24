package redispkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	args := m.Called(ctx, key, expiration)
	return args.Get(0).(*redis.BoolCmd)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, values)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	args := m.Called(ctx, key, values)
	return args.Get(0).(*redis.IntCmd)
}
