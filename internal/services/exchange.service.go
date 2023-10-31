package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"momentum-go-server/internal/models"
	"net/http"
	"time"
)

var currentTime = time.Now()
var yesterday = currentTime.AddDate(0, 0, -1)

func getNBUData(date string) []models.NBU {
	var response []models.NBU

	apiURL := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?date=%s&json", date)

	resp, err := http.Get(apiURL)

	if err != nil {
		log.Println("Error creating HTTP request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading HTTP response body:", err)
	}
	json.Unmarshal([]byte(body), &response)

	return response
}

func GetNBUExchange() models.ExchangeRatesResponse {
	frontendResponse := make(map[string]models.ExchangeFrontendResponse)

	yyyymmddNoDash := currentTime.Format("20060102")
	yyyymmddNoDashPreviousDay := yesterday.Format("20060102")

	todayData := getNBUData(yyyymmddNoDash)
	yesterdayData := getNBUData(yyyymmddNoDashPreviousDay)

	for _, nbu := range todayData {
		for _, rate := range yesterdayData {
			if nbu.Symbol == rate.Symbol {

				frontendResponse[nbu.Symbol] = models.ExchangeFrontendResponse{
					Change:  rate.Rate - nbu.Rate,
					EndRate: nbu.Rate,
				}
			}
		}
	}
	NBU := models.ExchangeRatesResponse{
		NBU: frontendResponse,
	}
	return NBU
}
