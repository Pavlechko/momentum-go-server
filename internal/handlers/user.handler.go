package handlers

import (
	"net/http"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !IsDecodeJSONRequest(w, r, &user) {
		return
	}

	token, err := services.CreateUser(user)

	if err != nil {
		err = WriteJSONError(w, http.StatusConflict, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
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
		err = WriteJSONError(w, http.StatusUnauthorized, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	}

	WriteToken(w, http.StatusOK, token)
}
