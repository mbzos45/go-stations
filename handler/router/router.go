package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

func NewRouter(todoDB *sql.DB) *negroni.Negroni {
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
	n := negroni.New()
	n.Use(negroni.HandlerFunc(middleware.SetClientOS))
	n.UseHandler(mux)
	return n
}
