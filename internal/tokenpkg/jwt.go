// Package tokenpkg is responsible for generating and validating JWT tokens.
package tokenpkg

import (
	"github.com/dgrijalva/jwt-go"
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"time"
)

// Claims is a struct that will be encoded to a JWT.
type Claims struct {
	IP           string `json:"ip"`
	MaxReqPerSec int    `json:"max_req_per_sec"`
	jwt.StandardClaims
}

// NewJWT generates a new JWT tokenpkg string.
// The tokenpkg will expire after the specified duration.
// The tokenpkg will contain the IP and the maximum number of requests per second.
func NewJWT(ip string, expirationDuration time.Duration, maxReqPerSec int) (string, error) {
	expirationTime := time.Now().Add(expirationDuration)
	claims := &Claims{
		IP:           ip,
		MaxReqPerSec: maxReqPerSec,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JwtKey returns the JWT key.
func JwtKey() []byte {
	return []byte(confpkg.Config.JWTKey)
}
