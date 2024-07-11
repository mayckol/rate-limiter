package middlewarepkg

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mayckol/rate-limiter/internal/tokenpkg"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("API_KEY")
		if authHeader == "" {
			http.Error(w, "API_KEY header is required", http.StatusUnauthorized)
			return
		}

		claims := &tokenpkg.Claims{}
		token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
			return tokenpkg.JwtKey(), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid tokenpkg", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// analize the number of request per ip by interval, maybe using a redis, then add
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
