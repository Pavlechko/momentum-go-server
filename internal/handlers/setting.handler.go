package handlers

import (
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func GetSettings(w http.ResponseWriter, r *http.Request) {

}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	utils.InfoLogger.Println("UpdateSettings")
	userId := utils.GetUserId(r)
	vars := mux.Vars(r)
	settingsType := vars["type"]

	switch settingsType {
	case "quote":
		newQuote := services.QuoteUpdate(userId)
		WriteJSON(w, http.StatusOK, newQuote)
	case "background":
		newBackground := services.BackgroundUpdate(userId)
		WriteJSON(w, http.StatusOK, newBackground)
	}
}
