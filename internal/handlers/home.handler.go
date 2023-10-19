package handlers

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/services"
	"net/http"
)

type WeatherData struct {
	OpenWeather     models.FrontendWeatherResponse
	TomorrowWeather models.FrontendWeatherResponse
}

type ExchangeData struct {
	Name string
}

type ResponseObj struct {
	Weather WeatherData
}

func Home(w http.ResponseWriter, r *http.Request) {

	openWeatherRes := services.GetOpenWeatherData()
	tomorrowWeatherRes := services.GetTomorrowWeatherData()

	Weather := WeatherData{
		OpenWeather:     openWeatherRes,
		TomorrowWeather: tomorrowWeatherRes,
	}

	Response := ResponseObj{
		Weather: Weather,
	}

	WriteJSON(w, http.StatusOK, Response)
}
