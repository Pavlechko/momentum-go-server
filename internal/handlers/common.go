package handlers

import (
	"net/http"

	"momentum-go-server/internal/services"
)

type Handler struct {
	Service services.IService
}

type IHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
	SignInHandler(w http.ResponseWriter, r *http.Request)
	SignUpHandler(w http.ResponseWriter, r *http.Request)
	UpdateSettings(w http.ResponseWriter, r *http.Request)
}

func CreateHandler(service services.IService) *Handler {
	return &Handler{Service: service}
}
