package main

import (
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/infra/httppkg/webserver"
	"log"
)

func main() {

	// Load the configuration
	appEnv := confpkg.GetEnvPath()
	/*
	   To handle not set env variables, just check the len of the invalidVars slice (second return value)
	*/
	_, _, err := confpkg.LoadConfig(appEnv)
	if err != nil {
		log.Fatalln(err)
	}

	webserver.Start()
}
