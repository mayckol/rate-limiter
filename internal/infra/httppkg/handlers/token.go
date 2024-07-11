package handlers

import (
	"fmt"
	"github.com/mayckol/rate-limiter/internal/tokenpkg"
	"net/http"
	"time"
)

func Token(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr
	token, err := tokenpkg.NewJWT(ip, 10*time.Second, 10)
	if err != nil {
		http.Error(w, "error generating tokenpkg", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Api-Key", fmt.Sprintf(token))
	w.Write([]byte("token generated"))
}
