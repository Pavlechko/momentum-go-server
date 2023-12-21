package services

import (
	"github.com/google/uuid"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func (s *Service) QuoteUpdate(userID string) models.QuoteResponse {
	newQuote := GetRandomQuote()

	id, _ := uuid.Parse(userID)
	v := map[string]string{
		"author":  newQuote.Author,
		"content": newQuote.Content,
	}
	res, err := s.DB.UpdateSetting(id, models.Quote, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating quote settings:", err)
		return newQuote
	}
	return models.QuoteResponse{
		Content: res.Value["content"],
		Author:  res.Value["author"],
	}
}

func (s *Service) BackgroundUpdate(userID, source string) models.FrontendBackgroundImageResponse {
	newBackground := GetRandomBackground(source)
	id, _ := uuid.Parse(userID)
	v := map[string]string{
		"photographer": newBackground.Photographer,
		"image":        newBackground.Image,
		"alt":          newBackground.Alt,
		"source":       newBackground.Source,
		"source_url":   newBackground.SourceURL,
	}
	res, err := s.DB.UpdateSetting(id, models.Background, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating background settings:", err)
		return newBackground
	}

	return models.FrontendBackgroundImageResponse{
		Photographer: res.Value["photographer"],
		Image:        res.Value["image"],
		Alt:          res.Value["alt"],
		Source:       res.Value["source"],
		SourceURL:    res.Value["source_url"],
	}
}

func (s *Service) WeatherUpdate(userID, source, city string) models.FrontendWeatherResponse {
	newWeather := GetNewWeatherData(source, city)
	id, _ := uuid.Parse(userID)

	v := map[string]string{
		"city":   city,
		"source": source,
	}
	_, err := s.DB.UpdateSetting(id, models.Weather, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating weather settings:", err)
		return newWeather
	}

	return newWeather
}

func (s *Service) ExchangeUpdate(userID, source, from, to string) models.ExchangeFrontendResponse {
	newExchange := GetNewExchange(source, from, to)
	id, _ := uuid.Parse(userID)
	v := map[string]string{
		"from":   from,
		"source": source,
		"to":     to,
	}
	_, err := s.DB.UpdateSetting(id, models.Exchange, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating exchange settings:", err)
		return newExchange
	}
	return newExchange
}

func (s *Service) MarketUpdate(userID, symbol string) models.StockMarketResponse {
	newMarket := s.GetMarket(symbol)
	id, _ := uuid.Parse(userID)
	v := map[string]string{
		"symbol": symbol,
	}
	_, err := s.DB.UpdateSetting(id, models.Market, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating market settings:", err)
		return newMarket
	}
	return newMarket
}

func (s *Service) GetSettingData(userID string) models.SettingResponse {
	id, _ := uuid.Parse(userID)
	settings := s.DB.GetSettings(id)
	return settings
}
