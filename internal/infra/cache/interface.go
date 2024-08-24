package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Close() error
}
