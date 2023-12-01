package handlers

import (
	"net/http"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !h.IsDecodeJSONRequest(w, r, &user) {
		return
	}

	token, err := h.Service.CreateUser(user)

	if err != nil {
		err = h.WriteJSONError(w, http.StatusConflict, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	}

	h.WriteToken(w, http.StatusOK, token)
}

func (h *Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !h.IsDecodeJSONRequest(w, r, &user) {
		return
	}

	token, err := h.Service.GetUser(user)

	if err != nil {
		err = h.WriteJSONError(w, http.StatusUnauthorized, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	}

	h.WriteToken(w, http.StatusOK, token)
}
