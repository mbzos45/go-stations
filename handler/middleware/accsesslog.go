package middleware

import (
	"encoding/json"
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

func Logger(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
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
	bytes, err := json.Marshal(accLog)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes))
}
