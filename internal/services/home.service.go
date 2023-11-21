package services

import "momentum-go-server/internal/models"

func GetData(userID string) models.ResponseObj {
	var Response models.ResponseObj

	QuoteRes := GetQuote(userID)
	BackgroundRes := GetBackgroundData(userID)
	WeatherRes := GetWeatherData(userID)
	ExchangeRes := GetExchange(userID)
	MarketRes := GetMarketData(userID)
	SettingRes := GetSettingData(userID)

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
