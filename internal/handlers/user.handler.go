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

	token, err := services.CreateUser(user)

	if err != nil {
		WriteJSONError(w, http.StatusConflict, err.Error())
	}

	WriteToken(w, http.StatusOK, token)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !IsDecodeJSONRequest(w, r, &user) {
		return
	}

	token, err := services.GetUser(user)

	if err != nil {
		WriteJSONError(w, http.StatusUnauthorized, err.Error())
	}

	WriteToken(w, http.StatusOK, token)
}
