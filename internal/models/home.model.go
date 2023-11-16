package models

type ResponseObj struct {
	Weather    FrontendWeatherResponse
	Quote      QuoteResponse
	Background FrontendBackgroundImageResponse
	Exchange   ExchangeFrontendResponse
	Market     StockMarketResponse
	Settings   SettingResponse
}
