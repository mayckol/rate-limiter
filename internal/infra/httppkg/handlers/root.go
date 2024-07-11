package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mayckol/rate-limiter/internal/infra/httppkg/middlewarepkg"
	"net/http"
)

func Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middlewarepkg.SetEnvsMiddleware)
	//r.Use(RateLimitMiddleware)

	r.Get("/token", Token)
	r.With(middlewarepkg.RateLimitMiddleware).With(middlewarepkg.AuthMiddleware).Get("/private", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a private endpoint"))
	})
	return r
}
