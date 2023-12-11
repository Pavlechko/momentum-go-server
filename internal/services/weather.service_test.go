package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestGetWindDirection(t *testing.T) {
	testCases := []struct {
		windSpeed, rawDirect float64
		expectedResult       string
	}{
		{10, 0, "10m/s N"},
		{10, 22.6, "10m/s NE"},
		{10, 67.6, "10m/s E"},
		{10, 112.6, "10m/s ES"},
		{10, 157.6, "10m/s S"},
		{10, 202.6, "10m/s SW"},
		{10, 245.6, "10m/s W"},
		{10, 300, "10m/s NW"},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case #%d", i+1), func(t *testing.T) {
			result := getWindDirection(testCase.windSpeed, testCase.rawDirect)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func TestGetWeatherDescription(t *testing.T) {
	testCases := []struct {
		rainIntensity, snowIntensity float64
		cloudCover                   int
		expectedIcon, expectedMain   string
	}{
		{0, 0, 20, "01d", "Clear"},
		{0, 0, 26, "02d", "Mostly sunny"},
		{2.1, 0, 26, "09d", "Havy rain"},
	}
	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Case #%d", i+1), func(t *testing.T) {
			icon, main := getWeatherDescription(testCase.rainIntensity, testCase.snowIntensity, testCase.cloudCover)
			assert.Equal(t, testCase.expectedIcon, icon)
			assert.Equal(t, testCase.expectedMain, main)
		})
	}
}

func TestGetWeatherData(t *testing.T) {
	uid, _ := uuid.NewRandom()
	id := uid.String()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockData(ctrl)

	service := &Service{
		DB: mockStore,
	}
	t.Run("Case when the source of the Weather is OpenWeather API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		godotenv.Load()

		mockSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Weather",
			Value: map[string]string{
				"city":   "Kyiv",
				"source": "OpenWeather",
			},
		}

		var (
			city   = mockSettingRes.Value["city"]
			source = mockSettingRes.Value["source"]
		)

		fakeOpenWeatherResponse := `
			{
				"main": {
					"temp": 14.2,
					"feels_like": 12.6,
					"humidity": 63
				},
				"wind": {
					"speed": 1.6,
					"deg": 268
				},
				"weather": [
					{
						"main": "Clear",
						"icon": "01d"
					}
				]
			}
		`

		mockStore.EXPECT().GetSettingByName(uid, models.Weather).Return(*mockSettingRes, nil)

		url := fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=",
			city,
		)

		httpmock.RegisterResponder(
			"GET",
			url+os.Getenv("OPEN_WEATHER_API"),
			httpmock.NewStringResponder(200, fakeOpenWeatherResponse),
		)

		weather := service.GetWeatherData(id)

		assert.Equal(t, source, weather.Source)
		assert.Equal(t, 14.2, weather.Temp)
	})

	t.Run("Case when the source of the Weather is TomorrowWeather API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		godotenv.Load()

		mockSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Weather",
			Value: map[string]string{
				"city":   "Odesa",
				"source": "TomorrowWeather",
			},
		}

		var (
			city   = mockSettingRes.Value["city"]
			source = mockSettingRes.Value["source"]
		)

		fakeTomorrowWeatherResponse := `
			{
				"data": {
					"values": {
						"temperature": 14.6,
						"temperatureApparent": 11.98,
						"humidity": 59,
						"windSpeed": 11.2,
						"windDirection": 189,
						"cloudCover": 0,
						"rainIntensity": 0,
						"snowIntensity": 0.1
					}
				}
			}
		`

		mockStore.EXPECT().GetSettingByName(uid, models.Weather).Return(*mockSettingRes, nil)

		url := fmt.Sprintf(
			"https://api.tomorrow.io/v4/weather/realtime?location=%s&apikey=",
			city,
		)

		httpmock.RegisterResponder(
			"GET",
			url+os.Getenv("TOMORROW_WEATHER_API"),
			httpmock.NewStringResponder(200, fakeTomorrowWeatherResponse),
		)

		weather := service.GetWeatherData(id)

		assert.Equal(t, source, weather.Source)
		assert.Equal(t, 14.6, weather.Temp)
	})
}
