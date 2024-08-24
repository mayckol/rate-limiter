package tokenpkg

import (
	"github.com/golang-jwt/jwt/v5"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

var testIP = "127.0.0.1"

func TestNewJWTGenerates(t *testing.T) {
	_, _, err := confpkg.LoadConfig(true)
	t.Run("NewJWT generates a new JWT token string", func(t *testing.T) {
		if err != nil {
			log.Fatalln(err)
		}
		expirationDuration := 5 * time.Minute
		maxReqPerSec := 10

		tokenString, err := NewJWT(testIP, expirationDuration, maxReqPerSec)
		assert.NoError(t, err, "Expected no error")

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtKey(), nil
		})

		assert.NoError(t, err, "Expected no error")
		assert.True(t, token.Valid, "Expected token to be valid")

		if claims, ok := token.Claims.(*Claims); ok {
			assert.Equal(t, testIP, claims.IP, "IP should match")
			assert.Equal(t, maxReqPerSec, claims.MaxReqPerSec, "MaxReqPerSec should match")
		} else {
			t.Fatal("Expected claims to be of type *Claims")
		}
	})

	t.Run("NewJWT generates a new JWT token string with zero maxReqPerSec", func(t *testing.T) {
		expirationDuration := 5 * time.Minute
		maxReqPerSec := 0

		tokenString, err := NewJWT(testIP, expirationDuration, maxReqPerSec)
		assert.NoError(t, err, "Expected no error")

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return JwtKey(), nil
		})

		assert.NoError(t, err, "Expected no error")

		if claims, ok := token.Claims.(*Claims); ok {
			assert.Equal(t, maxReqPerSec, claims.MaxReqPerSec, "MaxReqPerSec should be 0")
		} else {
			t.Fatal("Expected claims to be of type *Claims")
		}
	})
}
