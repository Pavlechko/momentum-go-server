package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func GetMarket(symbol string) models.StockMarketResponse {
	var response models.StockMarket
	var frontendResponse models.StockMarketResponse

	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}

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

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("json.Unmarshal response", err)
		return frontendResponse
	}

	frontendResponse = models.StockMarketResponse{
		Symbol:        response.Market.Symbol,
		Price:         response.Market.Price,
		Change:        response.Market.Change,
		ChangePercent: response.Market.ChangePercent,
	}

	return frontendResponse
}

func (s *Service) GetMarketData(userID string) models.StockMarketResponse {
	var response models.StockMarketResponse
	id, _ := uuid.Parse(userID)

	res, err := s.DB.GetSettingByName(id, models.Market)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Market setting:", err)
		return response
	}
	var symbol = res.Value["symbol"]
	response = GetMarket(symbol)
	return response
}
