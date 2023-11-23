package services

import (
	"sync"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

type Res struct {
	Mu      sync.Mutex
	Counter int
	Quit    chan bool
}

func (r *Res) GetData(userID string) models.ResponseObj {
	Response := models.ResponseObj{}
	cQuote := make(chan models.QuoteResponse)
	cBackground := make(chan models.FrontendBackgroundImageResponse)
	cWeather := make(chan models.FrontendWeatherResponse)
	cExchange := make(chan models.ExchangeFrontendResponse)
	cMarket := make(chan models.StockMarketResponse)
	cSetting := make(chan models.SettingResponse)

	go r.getQuote(cQuote, userID)
	go r.getBackground(cBackground, userID)
	go r.getWeather(cWeather, userID)
	go r.getExchange(cExchange, userID)
	go r.getMarket(cMarket, userID)
	go r.getSetting(cSetting, userID)

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
		case <-r.Quit:
			return Response
		}
	}
}

func (r *Res) getQuote(c chan<- models.QuoteResponse, userID string) {
	quoteRes := GetQuote(userID)
	c <- quoteRes
	r.checkCounter()
}

func (r *Res) getBackground(c chan<- models.FrontendBackgroundImageResponse, userID string) {
	backgroundRes := GetBackgroundData(userID)
	c <- backgroundRes
	r.checkCounter()
}

func (r *Res) getWeather(c chan<- models.FrontendWeatherResponse, userID string) {
	weatherRes := GetWeatherData(userID)
	c <- weatherRes
	r.checkCounter()
}

func (r *Res) getExchange(c chan<- models.ExchangeFrontendResponse, userID string) {
	exchangeRes := GetExchange(userID)
	c <- exchangeRes
	r.checkCounter()
}

func (r *Res) getMarket(c chan<- models.StockMarketResponse, userID string) {
	marketRes := GetMarketData(userID)
	c <- marketRes
	r.checkCounter()
}

func (r *Res) getSetting(c chan<- models.SettingResponse, userID string) {
	settingRes := GetSettingData(userID)
	c <- settingRes
	r.checkCounter()
}

func (r *Res) checkCounter() {
	r.Mu.Lock()
	r.Counter++
	r.Mu.Unlock()
	if r.Counter >= 6 {
		r.Quit <- true
		return
	}
}
