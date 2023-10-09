package handlers

import (
	"encoding/json"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"net/http"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := store.CreateUser(user) // call service

	WriteJSON(w, http.StatusOK, result)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := store.GetUser(user)

	WriteJSON(w, http.StatusOK, result)
}
