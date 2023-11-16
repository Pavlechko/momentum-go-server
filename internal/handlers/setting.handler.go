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
		var backgroundInput models.BackgroundInput
		if !IsDecodeJSONRequest(w, r, &backgroundInput) {
			utils.ErrorLogger.Println("Error decoding background input")
			return
		}
		var source = backgroundInput.Source
		for _, provider := range models.BACKGROUND_PROVIDERS {
			if source == provider {
				newBackground := services.BackgroundUpdate(userId, source)
				WriteJSON(w, http.StatusOK, newBackground)
				return
			}
		}
		WriteJSONError(w, http.StatusBadRequest, "No such provider found")
	case "weather":
		var weatherInput models.WeatherInput
		if !IsDecodeJSONRequest(w, r, &weatherInput) {
			utils.ErrorLogger.Println("Error decoding weather input")
			return
		}
		var (
			source    = weatherInput.Source
			inputCity = weatherInput.City
		)
		for _, provider := range models.WEATHER_PROVIDERS {
			for _, city := range models.CITIES {
				if source == provider && inputCity == city {
					newWeather := services.WeatherUpdate(userId, source, inputCity)
					WriteJSON(w, http.StatusOK, newWeather)
					return
				}
			}
		}
		WriteJSONError(w, http.StatusBadRequest, "No such provider or city found")
	case "exchange":
		var exchangeInput models.ExchangeInput
		if !IsDecodeJSONRequest(w, r, &exchangeInput) {
			utils.ErrorLogger.Println("Error decoding exchange input")
			return
		}
		var (
			source = exchangeInput.Source
			from   = exchangeInput.From
			to     = exchangeInput.To
		)
		newExchange := services.ExchangeUpdate(userId, source, from, to)
		WriteJSON(w, http.StatusOK, newExchange)
	case "market":
		var marketInput models.MarketInput
		if !IsDecodeJSONRequest(w, r, &marketInput) {
			utils.ErrorLogger.Println("Error decoding market input")
			return
		}
		var (
			symbol = marketInput.Symbol
		)
		newMarket := services.MarketUpdate(userId, symbol)
		WriteJSON(w, http.StatusOK, newMarket)
	}
}
