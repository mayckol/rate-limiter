package confpkg

import (
	"fmt"
	"github.com/mayckol/envsnatch"
	"github.com/mayckol/rate-limiter/utils"
)

var Config *Conf

type Conf struct {
	AppEnv              string `env:"APP_ENV"`
	WSHost              string `env:"WS_HOST"`
	JWTKey              string `env:"JWT_KEY"`
	RedisHost           string `env:"REDIS_HOST"`
	RedisPort           string `env:"REDIS_PORT"`
	RedisCacheKey       string `env:"REDIS_CACHE_KEY"`
	DefaultMaxReqPerSec int    `env:"DEFAULT_MAX_REQ_PER_SEC"`
	TokenExpiresInSec   int    `env:"TOKEN_EXPIRES_IN_SEC"`
	TimeoutDuration     int    `env:"TIMEOUT_DURATION"`
}

// LoadConfig loads the configuration from the .env file or .env.test file and returns the configuration and the invalid variables
func LoadConfig(isTest ...bool) (*Conf, *[]envsnatch.UnmarshalingErr, error) {
	path := utils.RootPath()
	envType := ".env"
	if isTest != nil && isTest[0] == true {
		envType = ".env.test"
	}
	es, _ := envsnatch.NewEnvSnatch()
	es.AddPath(path)
	es.AddFileName(envType)

	var cfg Conf
	invalidVars, err := es.Unmarshal(&cfg)
	if invalidVars != nil {
		for _, v := range *invalidVars {
			fmt.Printf("invalid var: %s reason: %s\n", v.Field, v.Reason)
		}
		return nil, invalidVars, err
	}

	Config = &cfg
	return &cfg, invalidVars, nil
}
