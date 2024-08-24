package webserver

import (
	"context"
	"errors"
	"fmt"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/entity"
	"github.com/mayckol/rate-limiter/internal/infra/httppkg/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(reqRepository entity.RequestRepositoryInterface) {
	server := &http.Server{Addr: confpkg.Config.WSHost, Handler: handlers.Handler(reqRepository)}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	fmt.Printf("starting server on %s\n", confpkg.Config.WSHost)
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
