package confpkg

import (
	"fmt"
	"github.com/mayckol/envsnatch"
	"os"
	"strings"
)

var Config *Conf

type Conf struct {
	WSHost string `env:"WS_HOST"`
	JWTKey string `env:"JWT_KEY"`
}

// LoadConfig loads the configuration from the .env file or .env.test file and returns the configuration and the invalid variables
func LoadConfig(envPath string) (*Conf, *[]envsnatch.UnmarshalingErr, error) {
	if envPath == "" {
		return nil, nil, fmt.Errorf("env path is empty")
	}
	file := strings.Join([]string{envPath}, "/")
	es, _ := envsnatch.NewEnvSnatch()
	es.AddPath(file[0:strings.LastIndex(file, "/")])
	es.AddFileName(file[strings.LastIndex(file, "/")+1:])

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

// GetEnvPath returns the path of the .env file, if isTestEnv is true, it returns the path of the .env.test file
func GetEnvPath(isTestEnv ...bool) string {
	absPath := os.Getenv("PWD")
	if len(isTestEnv) > 0 && isTestEnv[0] == true {
		return absPath + "/.env.test"
	}
	return absPath + "/.env"
}
