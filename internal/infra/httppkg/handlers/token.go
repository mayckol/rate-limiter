package handlers

import (
	confpkg "github.com/mayckol/rate-limiter/configpkg"
	"github.com/mayckol/rate-limiter/internal/tokenpkg"
	"net/http"
	"strconv"
	"time"
)

func Token(w http.ResponseWriter, r *http.Request) {
	maxReqPerSec := confpkg.Config.DefaultMaxReqPerSec
	queryReqPerSec := r.URL.Query().Get("max_req_per_sec")
	if queryReqPerSec != "" {
		reqPerSec, err := strconv.Atoi(queryReqPerSec)
		if err != nil || reqPerSec <= 0 {
			http.Error(w, "invalid max_req_per_sec", http.StatusBadRequest)
			return
		}
		maxReqPerSec = reqPerSec
	}

	tokenExpiresIn := time.Duration(confpkg.Config.TokenExpiresInSec) * time.Second
	expiration := r.URL.Query().Get("token_expires_in_sec")
	if expiration != "" {
		expires, err := strconv.Atoi(expiration)
		if err != nil || expires <= 0 {
			http.Error(w, "invalid token_expires_in_sec", http.StatusBadRequest)
			return
		}
		tokenExpiresIn = time.Duration(expires) * time.Second
	}

	ip := r.RemoteAddr
	token, err := tokenpkg.NewJWT(ip, tokenExpiresIn, maxReqPerSec)
	if err != nil {
		http.Error(w, "error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Api-Key", token)
	w.Write([]byte("token generated"))
}
