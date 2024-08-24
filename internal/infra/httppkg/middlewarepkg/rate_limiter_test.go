package middlewarepkg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/tokenpkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRequestRepository struct {
	mock.Mock
}

func (m *MockRequestRepository) CheckRateLimit(key string, maxReqPerSec int) (bool, error) {
	args := m.Called(key, maxReqPerSec)
	return args.Bool(0), args.Error(1)
}

func validToken() string {
	claims := &tokenpkg.Claims{
		IP:           "127.0.0.1",
		MaxReqPerSec: 10,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(tokenpkg.JwtKey())
	return tokenString
}

func expiredToken() string {
	claims := &tokenpkg.Claims{
		IP:           "127.0.0.1",
		MaxReqPerSec: 10,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(tokenpkg.JwtKey())
	return tokenString
}

func TestNewRateLimiterMiddleware(t *testing.T) {
	mockRepo := new(MockRequestRepository)
	middleware := NewRateLimiterMiddleware(mockRepo)
	assert.NotNil(t, middleware)
}
func TestSetJWTClaimsMiddleware(t *testing.T) {
	confpkg.LoadConfig(true)
	t.Run("Valid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", validToken())
		rr := httptest.NewRecorder()

		middleware := &MiddlewarePkg{}
		handler := middleware.SetJWTClaimsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*tokenpkg.Claims)
			assert.True(t, ok)
			assert.Equal(t, 10, claims.MaxReqPerSec)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("No token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		middleware := &MiddlewarePkg{}
		handler := middleware.SetJWTClaimsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*tokenpkg.Claims)
			assert.True(t, ok)
			assert.Equal(t, confpkg.Config.DefaultMaxReqPerSec, claims.MaxReqPerSec)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "invalid_token")
		rr := httptest.NewRecorder()

		middleware := &MiddlewarePkg{}
		handler := middleware.SetJWTClaimsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Expired token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", expiredToken())
		rr := httptest.NewRecorder()

		middleware := &MiddlewarePkg{}
		handler := middleware.SetJWTClaimsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Rate limit middleware error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "claims_err", &tokenpkg.Claims{
			IP:           "",
			MaxReqPerSec: 10,
		})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockRepo := new(MockRequestRepository)
		mockRepo.On("CheckRateLimit", "rate_limiter_", 10).Return(false, assert.AnError)

		middleware := &MiddlewarePkg{ReqRepository: mockRepo}
		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Rate limit middleware allows request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "claims", &tokenpkg.Claims{
			IP:           "127.0.0.1",
			MaxReqPerSec: 10,
		})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockRepo := new(MockRequestRepository)
		mockRepo.On("CheckRateLimit", "rate_limiter_127001", 10).Return(true, nil)

		middleware := &MiddlewarePkg{ReqRepository: mockRepo}
		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Rate limit middleware blocks request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "claims", &tokenpkg.Claims{
			IP:           "127.0.0.1",
			MaxReqPerSec: 10,
		})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockRepo := new(MockRequestRepository)
		mockRepo.On("CheckRateLimit", "rate_limiter_127001", 10).Return(false, nil)

		middleware := &MiddlewarePkg{ReqRepository: mockRepo}
		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Rate limit middleware error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		ctx := context.WithValue(req.Context(), "claims", &tokenpkg.Claims{
			IP:           "127.0.0.1",
			MaxReqPerSec: 10,
		})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		mockRepo := new(MockRequestRepository)
		mockRepo.On("CheckRateLimit", "rate_limiter_127001", 10).Return(false, assert.AnError)

		middleware := &MiddlewarePkg{ReqRepository: mockRepo}
		handler := middleware.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockRepo.AssertExpectations(t)
	})
}
