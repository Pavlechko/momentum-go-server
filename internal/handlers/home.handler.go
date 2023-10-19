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
	Weather models.WeatherData
}

func Home(w http.ResponseWriter, r *http.Request) {

	WeatherRes := services.GetWeatherData()

	Response := ResponseObj{
		Weather: WeatherRes,
	}

	WriteJSON(w, http.StatusOK, Response)
}
