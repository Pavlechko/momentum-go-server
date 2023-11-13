package services

import (
	"encoding/json"
	"fmt"
	"io"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var currentTime = time.Now()
var yesterday = currentTime.AddDate(0, 0, -1)

// var symbolsArr = []string{"AUD", "BRL", "EGP", "CAD", "CLP", "CNY", "CZK", "EGP", "EUR", "GBP", "HKD", "INR", "JPY", "KRW", "LTL", "LVL", "TRY", "USD", "XAG", "XAU", "UAH", "PLN"}

func getNBUData(date, symbol string) models.NBU {
	var response []models.NBU

	apiURL := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json", symbol, date)

	resp, err := http.Get(apiURL)

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return response[0]
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return response[0]
	}

	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
	}
	return response[0]
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

	apiURL := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/fluctuation?base=%s&start_date=%s&end_date=%s&symbols=%s", from, yyyymmddNoDashPreviousDay, yyyymmddNoDash, to)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
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

	json.Unmarshal([]byte(body), &response)

	for _, rate := range response.Rates {
		var todayRate float64
		if rate.EndRate != 0 {
			todayRate = 1 / rate.EndRate
		} else {
			todayRate = rate.EndRate
		}
		frontendResponse = models.ExchangeFrontendResponse{
			Change:  rate.Change,
			EndRate: todayRate,
			From:    from,
			To:      to,
			Source:  "Layer",
		}
	}

	return frontendResponse
}

func GetExchange() models.ExchangeFrontendResponse {
	// make req to DB
	// NBU := getNBUExchange("USD")
	Layer := getLayerExchange("UAH", "USD")

	return Layer
}
