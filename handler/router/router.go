package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

func NewRouter(todoDB *sql.DB) http.Handler {
	// register routes
	mux := http.NewServeMux()
	// assign Handler
	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)))
	mux.Handle("/do-panic", middleware.Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intended panic")
	})))
	mux.Handle("/client_os", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os, err := middleware.GetClientOS(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(os)
	}))
	mux.Handle("/auth", middleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		_, err := w.Write([]byte("authorized"))
		if err != nil {
			log.Println(err)
			return
		}
	})))
	mux.Handle("/sluggish_func", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/text")
		_, err := w.Write([]byte("Done sluggish_func"))
		if err != nil {
			log.Println(err)
			return
		}
	}))
	chain := alice.New(middleware.Recovery, middleware.SetUA, middleware.Logger)
	return chain.Then(mux)
}
