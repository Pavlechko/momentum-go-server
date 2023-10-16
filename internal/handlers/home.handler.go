package handlers

import (
	"net/http"
)

type WeatherData struct {
	Name string
}

type ExchangeData struct {
	Name string
}

type ResponseObj struct {
	Weather  WeatherData
	Exchange ExchangeData
}

var Weather = WeatherData{
	Name: "Weather",
}
var Ex = ExchangeData{
	Name: "Exchange",
}

func Home(w http.ResponseWriter, r *http.Request) {

	mockResponse := ResponseObj{
		Weather:  Weather,
		Exchange: Ex,
	}
	WriteJSON(w, http.StatusOK, mockResponse)
}
