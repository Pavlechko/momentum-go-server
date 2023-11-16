package services

import "momentum-go-server/internal/models"

func GetData(userId string) models.ResponseObj {
	var Response models.ResponseObj

	QuoteRes := GetQuote(userId)
	BackgroundRes := GetBackgroundData(userId)
	WeatherRes := GetWeatherData(userId)

	ExchangeRes := GetExchange()
	MarketRes := GetMarketData()
	SettingRes := GetSettingData(userId)

	Response = models.ResponseObj{
		Weather:    WeatherRes,
		Quote:      QuoteRes,
		Background: BackgroundRes,
		Exchange:   ExchangeRes,
		Market:     MarketRes,
		Settings:   SettingRes,
	}

	return Response
}
