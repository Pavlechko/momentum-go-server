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
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	json.Unmarshal([]byte(body), &response)

	return response
}
