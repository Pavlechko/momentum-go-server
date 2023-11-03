package services

import (
	"encoding/json"
	"io"
	"log"
	"momentum-go-server/internal/models"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetMarketData() models.StockMarketResponse {
	var response models.StockMarket

	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=DAX&apikey=" + os.Getenv("STOCK_MARKET_API"))

	if err != nil {
		log.Println("Error creating HTTP request:", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error reading HTTP response body:", err)
	}

	json.Unmarshal([]byte(body), &response)
	frontendResponse := models.StockMarketResponse{
		Symbol:        response.Market.Symbol,
		Price:         response.Market.Price,
		Change:        response.Market.Change,
		ChangePercent: response.Market.ChangePercent,
	}

	return frontendResponse
}
