package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"net/http"
)

type ExchangeData struct {
	Name string
}

type ResponseObj struct {
	Weather    models.WeatherData
	Quote      models.QuoteResponse
	Backgroung models.BackgroundData
	Exchange   []models.ExchangeResponse
}

func Home(w http.ResponseWriter, r *http.Request) {

	WeatherRes := services.GetWeatherData()
	QuoteRes := services.GetRandomQuote()
	BackgroundRes := services.GetBackgroundData()
	ExchangeRes := services.GetNBUExchange()

	Response := ResponseObj{
		Weather:    WeatherRes,
		Quote:      QuoteRes,
		Backgroung: BackgroundRes,
		Exchange:   ExchangeRes,
	}

	WriteJSON(w, http.StatusOK, Response)
}
