package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
)

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
		if slices.Contains(models.BACKGROUND_PROVIDERS, source) {
			newBackground := services.BackgroundUpdate(userId, source)
			WriteJSON(w, http.StatusOK, newBackground)
			return
		}
		WriteJSONError(w, http.StatusBadRequest, "No such provider found")
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
		if slices.Contains(models.WEATHER_PROVIDERS, source) && slices.Contains(models.CITIES, city) {
			newWeather := services.WeatherUpdate(userId, source, city)
			WriteJSON(w, http.StatusOK, newWeather)
			return
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

		if slices.Contains(models.EXCHANGE_PROVIDERS, source) && slices.Contains(models.CURRENCIES, from) && slices.Contains(models.CURRENCIES, to) {
			newExchange := services.ExchangeUpdate(userId, source, from, to)
			WriteJSON(w, http.StatusOK, newExchange)
			return
		}
		WriteJSONError(w, http.StatusBadRequest, "No such provider or currency found")
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
