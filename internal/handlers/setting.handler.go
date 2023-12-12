package handlers

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserID(r)
	vars := mux.Vars(r)
	settingsType := vars["type"]

	switch settingsType {
	case "quote":
		newQuote := h.Service.QuoteUpdate(userID)
		err := h.WriteJSON(w, http.StatusOK, newQuote)
		if err != nil {
			utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
		}
	case "background":
		var backgroundInput models.BackgroundInput
		if !h.IsDecodeJSONRequest(w, r, &backgroundInput) {
			utils.ErrorLogger.Println("Error decoding background input")
			return
		}
		var source = backgroundInput.Source
		if slices.Contains(models.BackgroundProviders, source) {
			newBackground := h.Service.BackgroundUpdate(userID, source)
			err := h.WriteJSON(w, http.StatusOK, newBackground)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := h.WriteJSONError(w, http.StatusBadRequest, "No such provider found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	case "weather":
		var weatherInput models.WeatherInput
		if !h.IsDecodeJSONRequest(w, r, &weatherInput) {
			utils.ErrorLogger.Println("Error decoding weather input")
			return
		}
		var (
			source = weatherInput.Source
			city   = weatherInput.City
		)
		if slices.Contains(models.WEATHER_PROVIDERS, source) && slices.Contains(models.CITIES, city) {
			newWeather := h.Service.WeatherUpdate(userID, source, city)
			err := h.WriteJSON(w, http.StatusOK, newWeather)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := h.WriteJSONError(w, http.StatusBadRequest, "No such provider or city found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	case "exchange":
		var exchangeInput models.ExchangeInput
		if !h.IsDecodeJSONRequest(w, r, &exchangeInput) {
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
			newExchange := h.Service.ExchangeUpdate(userID, source, from, to)
			err := h.WriteJSON(w, http.StatusOK, newExchange)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := h.WriteJSONError(w, http.StatusBadRequest, "No such provider or currency found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	case "market":
		var marketInput models.MarketInput
		if !h.IsDecodeJSONRequest(w, r, &marketInput) {
			utils.ErrorLogger.Println("Error decoding market input")
			return
		}
		var (
			symbol = marketInput.Symbol
		)
		if slices.Contains(models.COMPANIES, symbol) {
			newMarket := h.Service.MarketUpdate(userID, symbol)
			err := h.WriteJSON(w, http.StatusOK, newMarket)
			if err != nil {
				utils.ErrorLogger.Printf("Error write JSON %s", err.Error())
			}
			return
		}
		err := h.WriteJSONError(w, http.StatusBadRequest, "No such company found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	default:
		err := h.WriteJSONError(w, http.StatusNotFound, "No such setting found")
		if err != nil {
			utils.ErrorLogger.Printf("Error write errorJSON %s", err.Error())
		}
	}
}
