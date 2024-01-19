package handler

import (
	"encoding/json"
	"github.com/TechBowl-japan/go-stations/model"
	"log"
	"net/http"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthzResponse := &model.HealthzResponse{Message: "OK"}
	err := json.NewEncoder(w).Encode(healthzResponse)
	if err != nil {
		log.Print(err)
	}
}
