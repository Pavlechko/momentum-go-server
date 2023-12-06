package handlers

import (
	"fmt"
	"net/http"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

const minNameLeng = 3
const minPassLeng = 6

func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if !h.IsDecodeJSONRequest(w, r, &user) {
		return
	}

	err := validateUser(user)

	if err != nil {
		err = h.WriteJSONError(w, http.StatusConflict, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
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

	err := validateUser(user)

	if err != nil {
		err = h.WriteJSONError(w, http.StatusConflict, err.Error())
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
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

func validateUser(user models.UserInput) error {
	if len(user.Name) < minNameLeng {
		return fmt.Errorf("your name is too short, min 3 symbols")
	} else if len(user.Password) < minPassLeng {
		return fmt.Errorf("your password is too short, min 6 symbols")
	}

	return nil
}
