package services

import (
	"encoding/json"
	"fmt"
	"io"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

func GetMarket(symbol string) models.StockMarketResponse {
	var response models.StockMarket
	var frontendResponse models.StockMarketResponse

	apiURL := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=", symbol)

	resp, err := http.Get(apiURL + os.Getenv("STOCK_MARKET_API"))

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

func GetMarketData(userId string) models.StockMarketResponse {
	var response models.StockMarketResponse
	id, _ := uuid.Parse(userId)

	res, err := store.GetSettingByName(id, models.Market)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Market setting:", err)
		return response
	}
	var symbol = res.Value["symbol"]
	response = GetMarket(symbol)
	return response
}
