package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !IsDecodeJSONRequest(w, r, &user) {
		return
	}

	result, err := services.CreateUser(user)

	if err != nil {
		WriteJSON(w, http.StatusConflict, err.Error())
	}

	WriteJSON(w, http.StatusOK, result)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !IsDecodeJSONRequest(w, r, &user) {
		return
	}

	result, err := services.GetUser(user)

	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, err.Error())
	}

	WriteJSON(w, http.StatusOK, result)
}
