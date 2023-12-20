package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"momentum-go-server/test/mocks"
)

func TestUpdateSettings(t *testing.T) {
	ctrl := gomock.NewController(t)

	uid, _ := uuid.NewRandom()
	id := uid.String()

	mockService := mocks.NewMockIService(ctrl)
	handler := &Handler{
		Service: mockService,
	}

	user := &models.UserResponse{
		ID:   uid,
		Name: "TestUser",
	}

	fakeToken, err := utils.GenerateJWT(user)
	assert.NoError(t, err)

	router := mux.NewRouter()
	router.HandleFunc("/setting/{type}", handler.UpdateSettings).Methods("PUT")

	w := httptest.NewRecorder()

	t.Run("Quote", func(t *testing.T) {
		quoteRes := &models.QuoteResponse{}

		req := httptest.NewRequest("PUT", "/setting/quote", http.NoBody)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		mockService.EXPECT().QuoteUpdate(id).Return(*quoteRes).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Background", func(t *testing.T) {
		backgroundRes := &models.FrontendBackgroundImageResponse{}
		mockBackgroundInput := &models.BackgroundInput{
			Source: "unsplash.com",
		}
		source := mockBackgroundInput.Source

		body, err := json.Marshal(mockBackgroundInput)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/setting/background", bytes.NewBuffer(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		mockService.EXPECT().BackgroundUpdate(id, source).Return(*backgroundRes).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Weather", func(t *testing.T) {
		weatherRes := &models.FrontendWeatherResponse{}
		mockWeatherInput := &models.WeatherInput{
			Source: "OpenWeather",
			City:   "Kyiv",
		}
		source := mockWeatherInput.Source
		city := mockWeatherInput.City

		body, err := json.Marshal(mockWeatherInput)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/setting/weather", bytes.NewBuffer(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		mockService.EXPECT().WeatherUpdate(id, source, city).Return(*weatherRes).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Exchange", func(t *testing.T) {
		exchangeRes := &models.ExchangeFrontendResponse{}
		mockExchangeInput := &models.ExchangeInput{
			Source: "NBU",
			From:   "UAH",
			To:     "USD",
		}
		source := mockExchangeInput.Source
		from := mockExchangeInput.From
		to := mockExchangeInput.To

		body, err := json.Marshal(mockExchangeInput)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/setting/exchange", bytes.NewBuffer(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		mockService.EXPECT().ExchangeUpdate(id, source, from, to).Return(*exchangeRes).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Market", func(t *testing.T) {
		marketRes := &models.StockMarketResponse{}
		mockMarketInput := &models.MarketInput{
			Symbol: "BMWYY",
		}
		symbol := mockMarketInput.Symbol

		body, err := json.Marshal(mockMarketInput)
		assert.NoError(t, err)

		req := httptest.NewRequest("PUT", "/setting/market", bytes.NewBuffer(body))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		mockService.EXPECT().MarketUpdate(id, symbol).Return(*marketRes).Times(1)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid route variable", func(t *testing.T) {
		w := httptest.NewRecorder()

		req := httptest.NewRequest("PUT", "/setting/wrong", http.NoBody)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		var responseBody map[string]string

		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, "No such setting found", responseBody["Error"])
	})
}
