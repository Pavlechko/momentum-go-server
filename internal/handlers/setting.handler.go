package handlers

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
)

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	vars := mux.Vars(r)
	settingsType := vars["type"]

	switch settingsType {
	case "quote":
		newQuote := services.QuoteUpdate(userID)
		err := WriteJSON(w, http.StatusOK, newQuote)
		if err != nil {
			utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
		}
	case "background":
		var backgroundInput models.BackgroundInput
		if !IsDecodeJSONRequest(w, r, &backgroundInput) {
			utils.ErrorLogger.Println("Error decoding background input")
			return
		}
		var source = backgroundInput.Source
		if slices.Contains(models.BackgroundProviders, source) {
			newBackground := services.BackgroundUpdate(userID, source)
			err := WriteJSON(w, http.StatusOK, newBackground)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := WriteJSONError(w, http.StatusBadRequest, "No such provider found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
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
			newWeather := services.WeatherUpdate(userID, source, city)
			err := WriteJSON(w, http.StatusOK, newWeather)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := WriteJSONError(w, http.StatusBadRequest, "No such provider or city found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
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

		if slices.Contains(models.ExchangeProviders, source) &&
			slices.Contains(models.CURRENCIES, from) &&
			slices.Contains(models.CURRENCIES, to) {
			newExchange := services.ExchangeUpdate(userID, source, from, to)
			err := WriteJSON(w, http.StatusOK, newExchange)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := WriteJSONError(w, http.StatusBadRequest, "No such provider or currency found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	case "market":
		var marketInput models.MarketInput
		if !IsDecodeJSONRequest(w, r, &marketInput) {
			utils.ErrorLogger.Println("Error decoding market input")
			return
		}
		var (
			symbol = marketInput.Symbol
		)
		if slices.Contains(models.COMPANIES, symbol) {
			newMarket := services.MarketUpdate(userID, symbol)
			err := WriteJSON(w, http.StatusOK, newMarket)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := WriteJSONError(w, http.StatusBadRequest, "No such company found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	default:
		err := WriteJSONError(w, http.StatusNotFound, "No such setting found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	}
}
