package handler

import (
	"context"
	"encoding/json"
	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
	"io"
	"log"
	"net/http"
	"strconv"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		req := &model.CreateTODORequest{}
		deq := json.NewDecoder(r.Body)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(r.Body)
		if err := deq.Decode(&req); err != nil {
			log.Println(err)
		}
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Create(r.Context(), req)
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodPut:
		req := &model.UpdateTODORequest{}
		deq := json.NewDecoder(r.Body)
		if err := deq.Decode(&req); err != nil {
			log.Println(err)
		}
		if req.ID == 0 || req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Update(r.Context(), req)
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodGet:
		req := &model.ReadTODORequest{}
		q := r.URL.Query()
		if idStr := q.Get("prev_id"); idStr != "" {
			prevId, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.PrevID = prevId
		}
		if sizeStr := q.Get("size"); sizeStr != "" {
			size, err := strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			req.Size = size
		} else {
			req.Size = 3
		}
		res, err := h.Read(r.Context(), req)
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		req := &model.DeleteTODORequest{}
		deq := json.NewDecoder(r.Body)
		if err := deq.Decode(req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Delete(r.Context(), req)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		return
	}

}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO_item.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODO_items.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	return &model.ReadTODOResponse{TODOs: todos}, err
}

// Update handles the endpoint that updates the TODO_item.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	return &model.DeleteTODOResponse{}, err
}
