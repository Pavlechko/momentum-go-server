package services

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"
)

func CreateService(db *store.Database) *Service {
	return &Service{
		DB:      db,
		Counter: 0,
		Quit:    make(chan bool),
	}
}

func (s *Service) GetData(userID string) models.ResponseObj {
	Response := models.ResponseObj{}
	cQuote := make(chan models.QuoteResponse)
	cBackground := make(chan models.FrontendBackgroundImageResponse)
	cWeather := make(chan models.FrontendWeatherResponse)
	cExchange := make(chan models.ExchangeFrontendResponse)
	cMarket := make(chan models.StockMarketResponse)
	cSetting := make(chan models.SettingResponse)

	go s.getQuote(cQuote, userID)
	go s.getBackground(cBackground, userID)
	go s.getWeather(cWeather, userID)
	go s.getExchange(cExchange, userID)
	go s.getMarket(cMarket, userID)
	go s.getSetting(cSetting, userID)

	for {
		select {
		case Response.Quote = <-cQuote:
			utils.InfoLogger.Println("Write Quote data")
		case Response.Background = <-cBackground:
			utils.InfoLogger.Println("Write Background data")
		case Response.Weather = <-cWeather:
			utils.InfoLogger.Println("Write Weather data")
		case Response.Exchange = <-cExchange:
			utils.InfoLogger.Println("Write Exchange data")
		case Response.Market = <-cMarket:
			utils.InfoLogger.Println("Write Market data")
		case Response.Settings = <-cSetting:
			utils.InfoLogger.Println("Write Settings data")
		case <-s.Quit:
			return Response
		}
	}
}

func (s *Service) getQuote(c chan<- models.QuoteResponse, userID string) {
	quoteRes := s.GetQuote(userID)
	c <- quoteRes
	s.checkCounter()
}

func (s *Service) getBackground(c chan<- models.FrontendBackgroundImageResponse, userID string) {
	backgroundRes := s.GetBackgroundData(userID)
	c <- backgroundRes
	s.checkCounter()
}

func (s *Service) getWeather(c chan<- models.FrontendWeatherResponse, userID string) {
	weatherRes := s.GetWeatherData(userID)
	c <- weatherRes
	s.checkCounter()
}

func (s *Service) getExchange(c chan<- models.ExchangeFrontendResponse, userID string) {
	exchangeRes := s.GetExchange(userID)
	c <- exchangeRes
	s.checkCounter()
}

func (s *Service) getMarket(c chan<- models.StockMarketResponse, userID string) {
	marketRes := s.GetMarketData(userID)
	c <- marketRes
	s.checkCounter()
}

func (s *Service) getSetting(c chan<- models.SettingResponse, userID string) {
	settingRes := s.GetSettingData(userID)
	c <- settingRes
	s.checkCounter()
}

func (s *Service) checkCounter() {
	s.Mu.Lock()
	utils.InfoLogger.Println("Counter:", s.Counter)
	s.Counter++
	s.Mu.Unlock()
	if s.Counter >= 6 {
		s.Quit <- true
		s.Counter = 0
		return
	}
}
