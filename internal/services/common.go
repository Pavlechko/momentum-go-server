package services

import (
	"sync"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
)

type Service struct {
	DB      store.Data
	Mu      sync.Mutex
	Counter int
	Quit    chan bool
}

type IService interface {
	GetData(userID string) models.ResponseObj
	CreateUser(user models.UserInput) (string, error)
	GetUser(user models.UserInput) (string, error)
	BackgroundUpdate(userID string, source string) models.FrontendBackgroundImageResponse
	MarketUpdate(userID string, symbol string) models.StockMarketResponse
	QuoteUpdate(userID string) models.QuoteResponse
	WeatherUpdate(userID string, source string, city string) models.FrontendWeatherResponse
	ExchangeUpdate(userID string, source string, from string, to string) models.ExchangeFrontendResponse
}

func CreateService(db store.Data) *Service {
	return &Service{
		DB:      db,
		Counter: 0,
		Quit:    make(chan bool),
	}
}
