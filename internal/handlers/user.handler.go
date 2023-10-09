package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !IsDecodeJSONRequest(w, r, &user) {
		return
	}

	result, err := store.CreateUser(user) // call service

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

	result, err := store.GetUser(user) // call service

	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, err.Error())
	}

	WriteJSON(w, http.StatusOK, result)
}
