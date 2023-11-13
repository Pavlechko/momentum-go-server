package services

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"

	"github.com/google/uuid"
)

func QuoteUpdate(userId string) models.QuoteResponse {
	newQuote := GetRandomQuote()

	id, _ := uuid.Parse(userId)
	v := map[string]string{
		"author":  newQuote.Author,
		"content": newQuote.Content,
	}
	res, err := store.UpdateSetting(id, models.Quote, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating quote settings:", err)
		return newQuote
	}
	return models.QuoteResponse{
		Content: res.Value["content"],
		Author:  res.Value["author"],
	}
}

func BackgroundUpdate(userId string) models.FrontendBackgroundImageResponse {
	newBackground := GetRandomBackground()
	id, _ := uuid.Parse(userId)
	v := map[string]string{
		"photographer": newBackground.Photographer,
		"image":        newBackground.Image,
		"alt":          newBackground.Alt,
		"source":       newBackground.Source,
		"source_url":   newBackground.SourceUrl,
	}
	res, err := store.UpdateSetting(id, models.Background, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating background settings:", err)
		return newBackground
	}

	return models.FrontendBackgroundImageResponse{
		Photographer: res.Value["photographer"],
		Image:        res.Value["image"],
		Alt:          res.Value["alt"],
		Source:       res.Value["source"],
		SourceUrl:    res.Value["source_url"],
	}
}

func WeatherUpdate(userId, source, city string) models.FrontendWeatherResponse {
	newWeather := GetNewWeatherData(source, city)
	id, _ := uuid.Parse(userId)

	v := map[string]string{
		"city":   city,
		"source": source,
	}
	_, err := store.UpdateSetting(id, models.Weather, v)
	if err != nil {
		utils.ErrorLogger.Println("Error updating weather settings:", err)
		return newWeather
	}

	return newWeather
}
