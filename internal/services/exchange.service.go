package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

var currentTime = time.Now()
var yesterday = currentTime.AddDate(0, 0, -1)

func getNBUData(date, symbol string) models.NBU {
	var response []models.NBU
	var NBURes models.NBU
	var currentRate float64

	apiURL := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json", symbol, date)

	resp, err := http.Get(apiURL)

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return NBURes
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return NBURes
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return NBURes
	}

	if len(response) == 0 {
		utils.ErrorLogger.Println("Rresponse is empty")
		return NBURes
	}

	if response[0].Rate != 0 {
		currentRate = 1 / response[0].Rate
	} else {
		currentRate = response[0].Rate
	}

	NBURes = models.NBU{
		Rate:   currentRate,
		Symbol: response[0].Symbol,
	}

	return NBURes
}

func getNBUExchange(symbol string) models.ExchangeFrontendResponse {
	var frontendResponse models.ExchangeFrontendResponse

	yyyymmddNoDash := currentTime.Format("20060102")
	yyyymmddNoDashPreviousDay := yesterday.Format("20060102")

	todayData := getNBUData(yyyymmddNoDash, symbol)
	yesterdayData := getNBUData(yyyymmddNoDashPreviousDay, symbol)

	frontendResponse = models.ExchangeFrontendResponse{
		Change:  yesterdayData.Rate - todayData.Rate,
		EndRate: todayData.Rate,
		From:    "UAH",
		To:      symbol,
		Source:  "NBU",
	}

	return frontendResponse
}

func getLayerExchange(from, to string) models.ExchangeFrontendResponse {
	var response models.LayerResponse
	var frontendResponse models.ExchangeFrontendResponse

	yyyymmddNoDash := currentTime.Format("2006-01-02")
	yyyymmddNoDashPreviousDay := yesterday.Format("2006-01-02")

	client := &http.Client{}

	apiURL := fmt.Sprintf(
		"https://api.apilayer.com/exchangerates_data/fluctuation?base=%s&start_date=%s&end_date=%s&symbols=%s",
		from,
		yyyymmddNoDashPreviousDay,
		yyyymmddNoDash,
		to,
	)

	req, err := http.NewRequest("GET", apiURL, http.NoBody)
	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}

	err = godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}

	req.Header.Set("apikey", os.Getenv("LAYER_EXCHANGE_API"))

	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLogger.Println("Error sending HTTP request:", err)
		return frontendResponse
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return frontendResponse
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return frontendResponse
	}

	for _, rate := range response.Rates {
		frontendResponse = models.ExchangeFrontendResponse{
			Change:  rate.Change,
			EndRate: rate.EndRate,
			From:    from,
			To:      to,
			Source:  "Layer",
		}
	}

	return frontendResponse
}

func (s *Service) GetExchange(userID string) models.ExchangeFrontendResponse {
	var response models.ExchangeFrontendResponse
	id, _ := uuid.Parse(userID)

	res, err := s.DB.GetSettingByName(id, models.Exchange)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Exchange setting:", err)
		return response
	}
	var (
		source = res.Value["source"]
		from   = res.Value["from"]
		to     = res.Value["to"]
	)
	response = GetNewExchange(source, from, to)

	return response
}

func GetNewExchange(source, from, to string) models.ExchangeFrontendResponse {
	var newExchange models.ExchangeFrontendResponse
	switch source {
	case "NBU":
		newExchange = getNBUExchange(to)
	case "Layer":
		newExchange = getLayerExchange(from, to)
	}

	return newExchange
}
