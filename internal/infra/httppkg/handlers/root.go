package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mayckol/rate-limiter/internal/entity"
	"github.com/mayckol/rate-limiter/internal/infra/httppkg/middlewarepkg"
	"net/http"
)

func Handler(reqRepository entity.RequestRepositoryInterface) http.Handler {
	r := chi.NewRouter()

	m := middlewarepkg.NewRateLimiterMiddleware(reqRepository)
	r.Use(middleware.Logger)

	r.Get("/token", Token)

	r.With(m.SetJWTClaimsMiddleware, m.RateLimitMiddleware).Get("/rate-limiter-active", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("success"))
		if err != nil {
			return
		}
	})
	return r
}
