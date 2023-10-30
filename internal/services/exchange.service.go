package services

import (
	"encoding/json"
	"io"
	"log"
	"momentum-go-server/internal/models"
	"net/http"
)

func convertToExchange(res []models.NBU) []models.ExchangeResponse {
	var exchangeData []models.ExchangeResponse

	for _, nbu := range res {
		exchange := models.ExchangeResponse(nbu)
		exchangeData = append(exchangeData, exchange)
	}

	return exchangeData
}

func GetNBUExchange() []models.ExchangeResponse {
	var response []models.NBU

	resp, err := http.Get("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?json")

	if err != nil {
		log.Println("Error creating HTTP request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading HTTP response body:", err)
	}
	json.Unmarshal([]byte(body), &response)

	frontendResponse := convertToExchange(response)
	return frontendResponse
}
