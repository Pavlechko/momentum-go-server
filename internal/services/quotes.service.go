package services

import (
	"encoding/json"
	"io"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"net/http"
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
	json.Unmarshal([]byte(body), &response)

	return response
}
