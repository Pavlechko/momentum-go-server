package services

import (
	"encoding/json"
	"io"
	"log"
	"momentum-go-server/internal/models"
	"net/http"
)

func GetRandomQuote() models.QuoteResponse {
	var response models.QuoteResponse

	resp, err := http.Get("https://api.quotable.io/random")

	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return response
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading HTTP response body:", err)
		return response
	}
	json.Unmarshal([]byte(body), &response)

	return response
}
