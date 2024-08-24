// Package tokenpkg is responsible for generating and validating JWT tokens.
package tokenpkg

import (
	"github.com/golang-jwt/jwt/v5"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"time"
)

// Claims is a struct that will be encoded to a JWT.
type Claims struct {
	IP           string `json:"ip"`
	MaxReqPerSec int    `json:"max_req_per_sec"`
	jwt.RegisteredClaims
}

// NewJWT generates a new JWT token string.
// The token will expire after the specified duration.
// The token will contain the IP and the maximum number of requests per second.
func NewJWT(ip string, expirationDuration time.Duration, maxReqPerSec int) (string, error) {
	expirationTime := time.Now().Add(expirationDuration)
	claims := &Claims{
		IP:           ip,
		MaxReqPerSec: maxReqPerSec,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(JwtKey())

	return tokenString, nil
}

func JwtKey() []byte {
	return []byte(confpkg.Config.JWTKey)
}
