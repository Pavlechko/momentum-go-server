package services

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"momentum-go-server/test/mocks"
)

func TestGetMarketData(t *testing.T) {
	uid, _ := uuid.NewRandom()
	id := uid.String()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}

	response := `{
		"Global Quote": {
			"01. symbol": "DEX",
			"05. price":  "359"
		}
	}`

	httpmock.RegisterResponder(
		"GET",
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=DEX&apikey="+os.Getenv("STOCK_MARKET_API"),
		httpmock.NewStringResponder(200, response),
	)

	mockStore := mocks.NewMockData(gomock.NewController(t))
	service := &Service{
		DB: mockStore,
	}

	mockMarketDB := &models.Setting{
		UserID: uid,
		Name:   "Market",
		Value: map[string]string{
			"symbol": "DEX",
		},
	}

	mockStore.EXPECT().GetSettingByName(uid, models.Market).Return(*mockMarketDB, nil)

	market := service.GetMarketData(id)

	assert.Equal(t, "359", market.Price)
}
