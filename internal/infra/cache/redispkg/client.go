package redispkg

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type ClientSettings struct {
	Host     string
	Port     string
	Password string
	AppEnv   string
}

func NewRedisClient(conf *ClientSettings) (*redis.Client, error) {
	host, port, password := conf.Host, conf.Port, conf.Password

	if host == "" || port == "" {
		return nil, fmt.Errorf("redis configuration error: REDIS_HOST or REDIS_PORT is not set")
	}

	opts := &redis.Options{
		Addr:       fmt.Sprintf("%s:%s", host, port),
		Password:   password,
		MaxRetries: 3,
	}

	if conf.AppEnv != "local" && password != "" {
		opts.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	}

	client := redis.NewClient(opts)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %v", err)
	}

	return client, nil
}
