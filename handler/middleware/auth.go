package middleware

import (
	"log"
	"net/http"
	"os"
)

const (
	uidKey = "BASIC_AUTH_USER_ID"
	pwKey  = "BASIC_AUTH_PASSWORD"
)

func Auth(h http.Handler) http.Handler {
	uid := os.Getenv(uidKey)
	pw := os.Getenv(pwKey)
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, sec, ok := r.BasicAuth()
		if id != uid || sec != pw || !ok {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized"))
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
