package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"momentum-go-server/internal/utils"
	"net/http"
)

type ResponseObj struct {
	Weather    models.FrontendWeatherResponse
	Quote      models.QuoteResponse
	Background models.FrontendBackgroundImageResponse
	Exchange   models.ExchangeFrontendResponse
	Market     models.StockMarketResponse
	Settings   models.SettingResponse
}

func Home(w http.ResponseWriter, r *http.Request) {
	userId := utils.GetUserId(r)

	WeatherRes := services.GetWeatherData()
	QuoteRes := services.GetRandomQuote()
	BackgroundRes := services.GetBackgroundData()
	ExchangeRes := services.GetExchange()
	MarketRes := services.GetMarketData()
	SettingRes := services.GetSettingData(userId)

	Response := ResponseObj{
		Weather:    WeatherRes,
		Quote:      QuoteRes,
		Background: BackgroundRes,
		Exchange:   ExchangeRes,
		Market:     MarketRes,
		Settings:   SettingRes,
	}

	WriteJSON(w, http.StatusOK, Response)
}
