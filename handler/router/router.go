package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler"
	"net/http"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()
	// Handler 設定
	healthzHandler := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	return mux
}
