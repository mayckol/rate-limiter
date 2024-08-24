package repository

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/infra/cache"
	"time"
)

type RequestRepository struct {
	CacheClient cache.ClientInterface
}

func NewRequestRepository(cacheClient cache.ClientInterface) *RequestRepository {
	return &RequestRepository{CacheClient: cacheClient}
}

// CheckRateLimit checks if the request is allowed under the rate limit.
func (r *RequestRepository) CheckRateLimit(key string, limit int) (bool, error) {
	ctx := context.Background()
	current, err := r.CacheClient.Get(ctx, key).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, err
	}

	if current < limit {
		_, err := r.CacheClient.Set(ctx, key, current+1, 1*time.Second).Result()
		if err != nil {
			return false, err
		}
		return true, nil
	}

	delay := time.Duration(confpkg.Config.TimeoutDuration) * time.Second
	_, err = r.CacheClient.Set(ctx, key, current+1, delay).Result()
	return false, err
}

func (r *RequestRepository) SetRateLimit(key string, limit int) error {
	res := r.CacheClient.Set(context.Background(), key, limit, time.Minute)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
