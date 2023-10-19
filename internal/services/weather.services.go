package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"momentum-go-server/internal/models"

	_ "github.com/joho/godotenv/autoload"
)

func GetWeatherData() models.WeatherData {

	openWeatherRes := getOpenWeatherData()
	tomorrowWeatherRes := getTomorrowWeatherData()

	Weather := models.WeatherData{
		OpenWeather:     openWeatherRes,
		TomorrowWeather: tomorrowWeatherRes,
	}

	return Weather
}

func getWindDirection(windSpeed, rawDirect float64) string {
	if rawDirect >= 342.5 && rawDirect <= 22.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s W"
	} else if rawDirect >= 22.6 && rawDirect <= 67.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s NE"
	} else if rawDirect >= 67.6 && rawDirect <= 112.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s E"
	} else if rawDirect >= 112.6 && rawDirect <= 157.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s ES"
	} else if rawDirect >= 157.6 && rawDirect <= 202.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s S"
	} else if rawDirect >= 202.6 && rawDirect <= 245.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s SW"
	} else if rawDirect >= 245.6 && rawDirect <= 297.5 {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s W"
	} else {
		return fmt.Sprintf("%.0f", windSpeed) + "m/s NW"
	}
}

func getOpenWeatherData() models.FrontendWeatherResponse {
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

	direction := getWindDirection(response.Wind.Speed, response.Wind.Direction)

	frontendResponse := models.FrontendWeatherResponse{
		Temp:       response.Main.Temp,
		FeelsLike:  response.Main.FeelsLike,
		Humidity:   response.Main.Humidity,
		WindSpeed:  direction,
		WeaterMain: response.Weather[0].Main,
		WeaterIcon: response.Weather[0].Icon,
		City:       response.Name + ", " + response.Sys.Country,
		Sourse:     "OpenWeatherAPI",
	}
	return frontendResponse
}

func getTomorrowWeatherData() models.FrontendWeatherResponse {
	var response models.TomorrowWeatherResponse
	var weaterIcon string
	var weaterMain string

	resp, err := http.Get("https://api.tomorrow.io/v4/weather/realtime?location=kyiv&apikey=" + os.Getenv("TOMORROW_WEATHER_API"))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &response)

	direction := getWindDirection(response.Data.Values.WindSpeed, response.Data.Values.WindDirection)

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
		Temp:       response.Data.Values.Temperature,
		FeelsLike:  response.Data.Values.TemperatureApparent,
		Humidity:   response.Data.Values.Humidity,
		WindSpeed:  direction,
		WeaterMain: weaterMain,
		WeaterIcon: weaterIcon,
		City:       response.Location.Name,
		Sourse:     "Tomorrow.io API",
	}
	return frontendResponse
}
