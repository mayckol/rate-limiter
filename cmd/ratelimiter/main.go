package main

import (
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/infra/cache/redispkg"
	"github.com/mayckol/rate-limiter/internal/infra/httppkg/webserver"
	"github.com/mayckol/rate-limiter/internal/infra/repository"
	"log"
)

func main() {
	conf, _, err := confpkg.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cacheClient, err := redispkg.NewRedisClient(&redispkg.ClientSettings{
		Host:     conf.RedisHost,
		Port:     conf.RedisPort,
		Password: conf.RedisCacheKey,
		AppEnv:   conf.AppEnv,
	})

	// sample of how to use memcached or other cache client
	//cacheClient, err := memcached.NewMemCachedClient(&memcached.ClientSettings{})

	if err != nil {
		log.Fatalln(err)
	}

	requestRepository := repository.NewRequestRepository(cacheClient)

	webserver.Start(requestRepository)
}
