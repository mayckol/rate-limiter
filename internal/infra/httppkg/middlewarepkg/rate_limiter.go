package middlewarepkg

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/entity"
	"github.com/mayckol/rate-limiter/utils"
	"net/http"
	"time"

	"github.com/mayckol/rate-limiter/internal/tokenpkg"
)

type MiddlewarePkg struct {
	ReqRepository entity.RequestRepositoryInterface
}

func NewRateLimiterMiddleware(reqRepository entity.RequestRepositoryInterface) *MiddlewarePkg {
	return &MiddlewarePkg{ReqRepository: reqRepository}
}

// SetJWTClaimsMiddleware extracts the JWT token from the API_KEY header and sets the claims in the request context.
func (m *MiddlewarePkg) SetJWTClaimsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("API_KEY")
		duration := time.Duration(confpkg.Config.TokenExpiresInSec) * time.Second
		claims := &tokenpkg.Claims{
			IP:           r.RemoteAddr,
			MaxReqPerSec: confpkg.Config.DefaultMaxReqPerSec,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			},
		}

		if authHeader == "" {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		tokenClaims := &tokenpkg.Claims{}
		t, err := jwt.ParseWithClaims(authHeader, tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return tokenpkg.JwtKey(), nil
		})
		if err != nil || !t.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims.MaxReqPerSec = tokenClaims.MaxReqPerSec
		claims.ExpiresAt = tokenClaims.ExpiresAt

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RateLimitMiddleware limits the number of requests per IP based on the maxReqPerSec in the JWT token.
// This example uses an in-memory rate limiter for simplicity.
func (m *MiddlewarePkg) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*tokenpkg.Claims)
		if !ok {
			http.Error(w, "unable to retrieve claims", http.StatusInternalServerError)
			return
		}

		maxReqPerSec := claims.MaxReqPerSec
		key := "rate_limiter_" + utils.ExtractNumbers(claims.IP)

		allowed, err := m.ReqRepository.CheckRateLimit(key, maxReqPerSec)
		if err != nil {
			http.Error(w, "rate limiting error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
