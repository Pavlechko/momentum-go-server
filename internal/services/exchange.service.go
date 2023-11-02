package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"momentum-go-server/internal/models"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var currentTime = time.Now()
var yesterday = currentTime.AddDate(0, 0, -1)
var symbolsArr = []string{"AUD", "BRL", "BTC", "CAD", "CLP", "CNY", "CZK", "EGP", "EUR", "GBP", "HKD", "INR", "JPY", "KRW", "LTL", "LVL", "TRY", "USD", "XAG", "XAU"}

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

func isExistItemInArray(arr []string, item string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == item {
			return true
		}
	}
	return false
}

func convertArrToString(arr []string) string {
	str := strings.Join(arr, ",")
	return str
}

func getNBUExchange() map[string]models.ExchangeFrontendResponse {
	frontendResponse := make(map[string]models.ExchangeFrontendResponse)

	yyyymmddNoDash := currentTime.Format("20060102")
	yyyymmddNoDashPreviousDay := yesterday.Format("20060102")

	todayData := getNBUData(yyyymmddNoDash)
	yesterdayData := getNBUData(yyyymmddNoDashPreviousDay)

	for _, nbu := range todayData {
		isExistRate := isExistItemInArray(symbolsArr, nbu.Symbol)
		for _, rate := range yesterdayData {
			if nbu.Symbol == rate.Symbol && isExistRate {
				frontendResponse[nbu.Symbol] = models.ExchangeFrontendResponse{
					Change:  rate.Rate - nbu.Rate,
					EndRate: nbu.Rate,
				}
			}
		}
	}

	return frontendResponse
}

func getLayerExchange() map[string]models.ExchangeFrontendResponse {
	var response models.LayerResponse
	frontendResponse := make(map[string]models.ExchangeFrontendResponse)

	yyyymmddNoDash := currentTime.Format("2006-01-02")
	yyyymmddNoDashPreviousDay := yesterday.Format("2006-01-02")
	symbols := convertArrToString(symbolsArr)

	client := &http.Client{}

	apiURL := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/fluctuation?base=UAH&start_date=%s&end_date=%s&symbols=%s", yyyymmddNoDashPreviousDay, yyyymmddNoDash, symbols)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
	}

	req.Header.Set("apikey", os.Getenv("LAYER_EXCHANGE_API"))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading HTTP response body:", err)
	}

	json.Unmarshal([]byte(body), &response)

	for key, rate := range response.Rates {
		var todayRate float64
		if rate.EndRate != 0 {
			todayRate = 1 / rate.EndRate
		} else {
			todayRate = rate.EndRate
		}
		frontendResponse[key] = models.ExchangeFrontendResponse{
			Change:  rate.Change,
			EndRate: todayRate,
		}
	}

	return frontendResponse
}

func GetExchange() models.ExchangeRatesResponse {
	NBU := getNBUExchange()
	Layer := getLayerExchange()

	Exchange := models.ExchangeRatesResponse{
		NBU:   NBU,
		Layer: Layer,
	}

	return Exchange
}
