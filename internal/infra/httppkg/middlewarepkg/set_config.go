package middlewarepkg

import (
	"context"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"net/http"
)

func SetEnvsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "conf", confpkg.Config)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
