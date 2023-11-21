package services

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"
)

func GetRandomQuote() models.QuoteResponse {
	var response models.QuoteResponse

	resp, err := http.Get("https://api.quotable.io/random")

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return response
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return response
	}

	return response
}

func GetQuote(userID string) models.QuoteResponse {
	const oneDayHours = 24
	var response models.QuoteResponse
	currentTime := time.Now()
	id, _ := uuid.Parse(userID)

	res, err := store.GetSettingByName(id, models.Quote)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Quote setting:", err)
		return response
	}
	dif := currentTime.Sub(res.UpdatedAt).Hours()

	if dif < oneDayHours {
		response = models.QuoteResponse{
			Content: res.Value["content"],
			Author:  res.Value["author"],
		}
		return response
	}
	return QuoteUpdate(userID)
}
