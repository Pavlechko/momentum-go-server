package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func (s *Service) GetWeatherData(userID string) models.FrontendWeatherResponse {
	var response models.FrontendWeatherResponse
	id, _ := uuid.Parse(userID)

	res, err := s.DB.GetSettingByName(id, models.Weather)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Weather setting:", err)
		return response
	}

	response = GetNewWeatherData(res.Value["source"], res.Value["city"])

	return response
}

func getWindDirection(windSpeed, rawDirect float64) string {
	const west = "m/s W"
	if rawDirect >= 342.5 && rawDirect <= 22.5 {
		return fmt.Sprintf("%.0f", windSpeed) + west
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
		return fmt.Sprintf("%.0f", windSpeed) + west
	}
	return fmt.Sprintf("%.0f", windSpeed) + "m/s NW"
}

func getOpenWeatherData(city string) models.FrontendWeatherResponse {
	var response models.OpenWeatherResponse
	var frontendResponse models.FrontendWeatherResponse

	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}

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
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return frontendResponse
	}

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

func getWeatherDescription(rainIntensity, snowIntensity float64, cloudCover int) (weaterIcon, weaterMain string) {
	const snowIcon = "13d"

	if rainIntensity <= 0.1 && snowIntensity <= 0.1 && (cloudCover >= 0 && cloudCover <= 25) {
		weaterIcon = "01d"
		weaterMain = "Clear"
		return weaterIcon, weaterMain
	} else if rainIntensity <= 0.1 && snowIntensity <= 0.1 && (cloudCover >= 26 && cloudCover <= 50) {
		weaterIcon = "02d"
		weaterMain = "Mostly sunny"
		return weaterIcon, weaterMain
	} else if rainIntensity <= 0.1 && snowIntensity <= 0.1 && (cloudCover >= 51 && cloudCover <= 75) {
		weaterIcon = "03d"
		weaterMain = "Mostly cloudy"
		return weaterIcon, weaterMain
	} else if rainIntensity <= 0.1 && snowIntensity <= 0.1 && (cloudCover >= 76 && cloudCover <= 100) {
		weaterIcon = "04d"
		weaterMain = "Clouds"
		return weaterIcon, weaterMain
	} else if rainIntensity >= 0.1 && rainIntensity <= 2 {
		weaterIcon = "10d"
		weaterMain = "Light rain"
		return weaterIcon, weaterMain
	} else if rainIntensity >= 2.1 && rainIntensity <= 6 {
		weaterIcon = "09d"
		weaterMain = "Havy rain"
		return weaterIcon, weaterMain
	} else if rainIntensity >= 6.1 {
		weaterIcon = "11d"
		weaterMain = "Thunderstorm"
		return weaterIcon, weaterMain
	} else if snowIntensity >= 0.1 && snowIntensity <= 2 {
		weaterIcon = snowIcon
		weaterMain = "Light snow"
		return weaterIcon, weaterMain
	} else if snowIntensity >= 2.1 && snowIntensity <= 6 {
		weaterIcon = snowIcon
		weaterMain = "Havy snow"
		return weaterIcon, weaterMain
	} else if snowIntensity >= 6.1 {
		weaterIcon = snowIcon
		weaterMain = "Blizzard"
		return weaterIcon, weaterMain
	}
	weaterIcon = "01d"
	weaterMain = "Clear"
	return weaterIcon, weaterMain
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
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return frontendResponse
	}

	direction := getWindDirection(response.Data.Values.WindSpeed, response.Data.Values.WindDirection)

	weaterIcon, weaterMain = getWeatherDescription(
		response.Data.Values.RainIntensity,
		response.Data.Values.SnowIntensity,
		response.Data.Values.CloudCover,
	)

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
