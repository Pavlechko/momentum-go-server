package services

import (
	"encoding/json"
	"io"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetMarketData() models.StockMarketResponse {
	var response models.StockMarket
	var frontendResponse models.StockMarketResponse

	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=DAX&apikey=" + os.Getenv("STOCK_MARKET_API"))

	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return frontendResponse
	}

	json.Unmarshal([]byte(body), &response)
	frontendResponse = models.StockMarketResponse{
		Symbol:        response.Market.Symbol,
		Price:         response.Market.Price,
		Change:        response.Market.Change,
		ChangePercent: response.Market.ChangePercent,
	}

	return frontendResponse
}
