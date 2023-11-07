package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"net/http"
)

type ResponseObj struct {
	Weather    models.WeatherData
	Quote      models.QuoteResponse
	Backgroung models.BackgroundData
	Exchange   models.ExchangeRatesResponse
	Market     models.StockMarketResponse
}

func Home(w http.ResponseWriter, r *http.Request) {

	WeatherRes := services.GetWeatherData()
	QuoteRes := services.GetRandomQuote()
	BackgroundRes := services.GetBackgroundData()
	ExchangeRes := services.GetExchange()
	MarketRes := services.GetMarketData()

	Response := ResponseObj{
		Weather:    WeatherRes,
		Quote:      QuoteRes,
		Backgroung: BackgroundRes,
		Exchange:   ExchangeRes,
		Market:     MarketRes,
	}

	WriteJSON(w, http.StatusOK, Response)
}
