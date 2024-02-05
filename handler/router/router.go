package router

import (
	"database/sql"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/urfave/negroni"
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
	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(mux)
	return n
}
