package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"momentum-go-server/internal/services"
)

type ErrorMessage struct {
	Error string
}

type Handler struct {
	Service *services.Service
}

func CreateHandler(service *services.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) WriteToken(w http.ResponseWriter, status int, token string) {
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.WriteHeader(status)
}

func (h *Handler) WriteJSONError(w http.ResponseWriter, status int, d string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(ErrorMessage{Error: d})
}

func (h *Handler) WriteJSON(w http.ResponseWriter, status int, d any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(d)
}

func (h *Handler) IsDecodeJSONRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return false
	}
	return true
}
