package redispkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("Set", mock.Anything, "key1", "value1", 10*time.Minute).Return(redis.NewStatusResult("OK", nil))
	mockRedis.Set(context.Background(), "key1", "value1", 10*time.Minute)
	mockRedis.AssertExpectations(t)
}

func TestGet(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("Get", mock.Anything, "key1").Return(redis.NewStringResult("value1", nil))
	mockRedis.Get(context.Background(), "key1")
	mockRedis.AssertExpectations(t)
}

func TestDel(t *testing.T) {
	mockRedis := new(MockRedisClient)

	cmd := redis.NewIntCmd(context.Background(), "DEL", "key1")
	cmd.SetVal(1)

	mockRedis.On("Del", mock.Anything, []string{"key1"}).Return(cmd)

	result := mockRedis.Del(context.Background(), "key1")

	if result.Val() != 1 {
		t.Errorf("Expected 1, got %d", result.Val())
	}

	mockRedis.AssertExpectations(t)
}

func TestRPush(t *testing.T) {
	mockRedis := new(MockRedisClient)

	cmd := redis.NewIntCmd(context.Background(), "RPUSH", "key1", "value1")
	cmd.SetVal(1)

	mockRedis.On("RPush", mock.Anything, "key1", []interface{}{"value1"}).Return(cmd)

	result := mockRedis.RPush(context.Background(), "key1", "value1")

	if result.Val() != 1 {
		t.Errorf("Expected 1, got %d", result.Val())
	}

	mockRedis.AssertExpectations(t)
}

func TestHSet(t *testing.T) {
	mockRedis := new(MockRedisClient)

	cmd := redis.NewIntCmd(context.Background(), "HSET", "key1", "field1", "value1")
	cmd.SetVal(1)

	mockRedis.On("HSet", mock.Anything, "key1", []interface{}{"field1", "value1"}).Return(cmd)

	result := mockRedis.HSet(context.Background(), "key1", "field1", "value1")

	if result.Val() != 1 {
		t.Errorf("Expected 1, got %d", result.Val())
	}

	mockRedis.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("Close").Return(nil)
	mockRedis.Close()
	mockRedis.AssertExpectations(t)
}
