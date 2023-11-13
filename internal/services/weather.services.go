package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"

	_ "github.com/joho/godotenv/autoload"
)

func GetWeatherData() models.FrontendWeatherResponse {
	// make request to DB
	openWeatherRes := getOpenWeatherData("kyiv")
	// tomorrowWeatherRes := getTomorrowWeatherData("kyiv")

	return openWeatherRes
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

func getOpenWeatherData(city string) models.FrontendWeatherResponse {
	var response models.OpenWeatherResponse
	var frontendResponse models.FrontendWeatherResponse

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=", city)

	resp, err := http.Get(url + os.Getenv("OPEN_WEATHER_API"))

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return frontendResponse
	}
	json.Unmarshal([]byte(body), &response)

	direction := getWindDirection(response.Wind.Speed, response.Wind.Direction)

	frontendResponse = models.FrontendWeatherResponse{
		Temp:       response.Main.Temp,
		FeelsLike:  response.Main.FeelsLike,
		Humidity:   response.Main.Humidity,
		WindSpeed:  direction,
		WeaterMain: response.Weather[0].Main,
		WeaterIcon: response.Weather[0].Icon,
		City:       city,
		Source:     "OpenWeather",
	}
	return frontendResponse
}

func getTomorrowWeatherData(city string) models.FrontendWeatherResponse {
	var response models.TomorrowWeatherResponse
	var frontendResponse models.FrontendWeatherResponse
	var weaterIcon string
	var weaterMain string

	url := fmt.Sprintf("https://api.tomorrow.io/v4/weather/realtime?location=%s&apikey=", city)

	resp, err := http.Get(url + os.Getenv("TOMORROW_WEATHER_API"))

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return frontendResponse
	}
	json.Unmarshal([]byte(body), &response)

	direction := getWindDirection(response.Data.Values.WindSpeed, response.Data.Values.WindDirection)

	if response.Data.Values.RainIntensity <= 0.1 && response.Data.Values.SnowIntensity <= 0.1 && (response.Data.Values.CloudCover >= 0 && response.Data.Values.CloudCover <= 25) {
		weaterIcon = "01d"
		weaterMain = "Clear"
	} else if response.Data.Values.RainIntensity <= 0.1 && response.Data.Values.SnowIntensity <= 0.1 && (response.Data.Values.CloudCover >= 26 && response.Data.Values.CloudCover <= 50) {
		weaterIcon = "02d"
		weaterMain = "Mostly sunny"
	} else if response.Data.Values.RainIntensity <= 0.1 && response.Data.Values.SnowIntensity <= 0.1 && (response.Data.Values.CloudCover >= 51 && response.Data.Values.CloudCover <= 75) {
		weaterIcon = "03d"
		weaterMain = "Mostly cloudy"
	} else if response.Data.Values.RainIntensity <= 0.1 && response.Data.Values.SnowIntensity <= 0.1 && (response.Data.Values.CloudCover >= 76 && response.Data.Values.CloudCover <= 100) {
		weaterIcon = "04d"
		weaterMain = "Clouds"
	} else if response.Data.Values.RainIntensity >= 0.1 && response.Data.Values.RainIntensity <= 2 {
		weaterIcon = "10d"
		weaterMain = "Light rain"
	} else if response.Data.Values.RainIntensity >= 2.1 && response.Data.Values.RainIntensity <= 6 {
		weaterIcon = "09d"
		weaterMain = "Havy rain"
	} else if response.Data.Values.RainIntensity >= 6.1 {
		weaterIcon = "11d"
		weaterMain = "Thunderstorm"
	} else if response.Data.Values.SnowIntensity >= 0.1 && response.Data.Values.SnowIntensity <= 2 {
		weaterIcon = "13d"
		weaterMain = "Light snow"
	} else if response.Data.Values.SnowIntensity >= 2.1 && response.Data.Values.SnowIntensity <= 6 {
		weaterIcon = "13d"
		weaterMain = "Havy snow"
	} else if response.Data.Values.SnowIntensity >= 6.1 {
		weaterIcon = "13d"
		weaterMain = "Blizzard"
	}

	frontendResponse = models.FrontendWeatherResponse{
		Temp:       response.Data.Values.Temperature,
		FeelsLike:  response.Data.Values.TemperatureApparent,
		Humidity:   response.Data.Values.Humidity,
		WindSpeed:  direction,
		WeaterMain: weaterMain,
		WeaterIcon: weaterIcon,
		City:       city,
		Source:     "TomorrowWeather",
	}
	return frontendResponse
}

func GetNewWeatherData(source, city string) models.FrontendWeatherResponse {
	var newWeather models.FrontendWeatherResponse
	switch source {
	case "OpenWeather":
		newWeather = getOpenWeatherData(city)
	case "TomorrowWeather":
		newWeather = getTomorrowWeatherData(city)
	}
	return newWeather
}
