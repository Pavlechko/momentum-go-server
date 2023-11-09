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
