package middleware

import (
	"log"
	"net/http"
	"time"
)

type AccessLog struct {
	Timestamp time.Time `json:"timestamp"`
	Latency   int64     `json:"latency"`
	Path      string    `json:"path"`
	OS        string    `json:"os"`
}

func Logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		os, err := GetClientOS(r.Context())
		if err != nil {
			log.Println(err)
			return
		}
		accLog := AccessLog{
			Timestamp: start,
			Latency:   time.Since(start).Milliseconds(),
			Path:      r.URL.Path,
			OS:        os,
		}
		log.Printf("%+v\n", accLog)
	}
	return http.HandlerFunc(fn)
}
