package services

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"momentum-go-server/internal/models"

	_ "github.com/joho/godotenv/autoload"
)

func GetOpenWeatherData() models.FrontendWeatherResponse {
	var response models.OpenWeatherResponse

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=kyiv&units=metric&appid=" + os.Getenv("OPEN_WEATHER_API"))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &response)

	frontendResponse := models.FrontendWeatherResponse{
		Temp:      response.Main.Temp,
		FeelsLike: response.Main.FeelsLike,
		Humidity:  response.Main.Humidity,
		WindSpeed: response.Wind.Speed,
		// WeaterDescription: response.Weather[0].Description,
		WeaterMain: response.Weather[0].Main,
		WeaterIcon: response.Weather[0].Icon,
		// Country:    response.Sys.Country,
		City:   response.Name + ", " + response.Sys.Country,
		Sourse: "OpenWeatherAPI",
	}
	// sb := string(body)
	return frontendResponse
}

func GetTomorrowWeatherData() models.FrontendWeatherResponse {
	var response models.TomorrowWeatherResponse
	var weaterIcon string
	var weaterMain string
	// var CloudCover string

	resp, err := http.Get("https://api.tomorrow.io/v4/weather/realtime?location=kyiv&apikey=" + os.Getenv("TOMORROW_WEATHER_API"))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &response)

	if response.Data.Values.RainIntensity == 0 && response.Data.Values.SnowIntensity == 0 && (response.Data.Values.CloudCover >= 0 && response.Data.Values.CloudCover <= 25) {
		weaterIcon = "01d"
		weaterMain = "Clear"
	} else if response.Data.Values.RainIntensity == 0 && response.Data.Values.SnowIntensity == 0 && (response.Data.Values.CloudCover >= 26 && response.Data.Values.CloudCover <= 50) {
		weaterIcon = "02d"
		weaterMain = "Mostly sunny"
	} else if response.Data.Values.RainIntensity == 0 && response.Data.Values.SnowIntensity == 0 && (response.Data.Values.CloudCover >= 51 && response.Data.Values.CloudCover <= 75) {
		weaterIcon = "03d"
		weaterMain = "Mostly cloudy"
	} else if response.Data.Values.RainIntensity == 0 && response.Data.Values.SnowIntensity == 0 && (response.Data.Values.CloudCover >= 76 && response.Data.Values.CloudCover <= 100) {
		weaterIcon = "04d"
		weaterMain = "Clouds"
	} else if response.Data.Values.RainIntensity == 1 || response.Data.Values.RainIntensity == 2 {
		weaterIcon = "10d"
		weaterMain = "Light rain"
	} else if response.Data.Values.RainIntensity >= 3 && response.Data.Values.RainIntensity <= 6 {
		weaterIcon = "09d"
		weaterMain = "Havy rain"
	} else if response.Data.Values.RainIntensity == 7 {
		weaterIcon = "11d"
		weaterMain = "Thunderstorm"
	} else if response.Data.Values.SnowIntensity == 1 || response.Data.Values.SnowIntensity == 2 {
		weaterIcon = "13d"
		weaterMain = "Light snow"
	} else if response.Data.Values.SnowIntensity >= 3 && response.Data.Values.SnowIntensity <= 6 {
		weaterIcon = "13d"
		weaterMain = "Havy snow"
	} else if response.Data.Values.SnowIntensity == 7 {
		weaterIcon = "13d"
		weaterMain = "Blizzard"
	}

	frontendResponse := models.FrontendWeatherResponse{
		Temp:      response.Data.Values.Temperature,
		FeelsLike: response.Data.Values.TemperatureApparent,
		Humidity:  response.Data.Values.Humidity,
		WindSpeed: response.Data.Values.WindSpeed,
		// WeaterDescription: response.Weather[0].Description,
		WeaterMain: weaterMain,
		WeaterIcon: weaterIcon,
		City:       response.Location.Name,
		Sourse:     "Tomorrow.io API",
	}
	// sb := string(body)
	return frontendResponse
}
