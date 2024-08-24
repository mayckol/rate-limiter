package memcached

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type ClientSettings struct {
}

func NewMemCachedClient(conf *ClientSettings) (interface{}, error) {
	return nil, nil
}

func (c ClientSettings) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	panic("implement me")
}

func (c ClientSettings) Get(ctx context.Context, key string) *redis.StringCmd {
	panic("implement me")
}

func (c ClientSettings) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	panic("implement me")
}

func (c ClientSettings) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (c ClientSettings) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	panic("implement me")
}

func (c ClientSettings) Incr(ctx context.Context, key string) *redis.IntCmd {
	panic("implement me")
}

func (c ClientSettings) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	panic("implement me")
}

func (c ClientSettings) Close() error {
	panic("implement me")
}
