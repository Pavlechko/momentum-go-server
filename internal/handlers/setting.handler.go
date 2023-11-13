package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func GetSettings(w http.ResponseWriter, r *http.Request) {
	// TO-DO
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
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
	case "weather":
		var weatherInput models.WeatherInput
		if !IsDecodeJSONRequest(w, r, &weatherInput) {
			utils.ErrorLogger.Println("Error decoding weather input")
			return
		}
		var (
			source = weatherInput.Source
			city   = weatherInput.City
		)

		newWeather := services.WeatherUpdate(userId, source, city)
		WriteJSON(w, http.StatusOK, newWeather)
	}
}
