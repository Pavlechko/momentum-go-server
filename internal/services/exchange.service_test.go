package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestGetExchange(t *testing.T) {
	uid, _ := uuid.NewRandom()
	id := uid.String()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var currentTime = time.Now()
	var yesterday = currentTime.AddDate(0, 0, -1)

	mockStore := mocks.NewMockData(ctrl)

	service := &Service{
		DB: mockStore,
	}

	t.Run("Case when the source of the Exchange is NBU API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Background",
			Value: map[string]string{
				"from":   "UAH",
				"to":     "USD",
				"source": "NBU",
			},
		}
		var (
			to                        = mockSettingRes.Value["to"]
			source                    = mockSettingRes.Value["source"]
			yyyymmddNoDash            = currentTime.Format("20060102")
			yyyymmddNoDashPreviousDay = yesterday.Format("20060102")
		)

		fakeNBUResponse := `[{
			"rate": 1,
			"cc": "USD"
		}]`

		fakePreviousNBUResponse := `[{
			"rate": 2.5,
			"cc": "USD"
		}]`

		mockStore.EXPECT().GetSettingByName(uid, models.Exchange).Return(*mockSettingRes, nil)

		httpmock.RegisterResponder(
			"GET",
			fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json", to, yyyymmddNoDash),
			httpmock.NewStringResponder(200, fakeNBUResponse),
		)

		httpmock.RegisterResponder(
			"GET",
			fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&date=%s&json", to, yyyymmddNoDashPreviousDay),
			httpmock.NewStringResponder(200, fakePreviousNBUResponse),
		)

		exchange := service.GetExchange(id)
		assert.Equal(t, source, exchange.Source)
		assert.Equal(t, -0.6, exchange.Change)
		assert.Equal(t, 1.0, exchange.EndRate)
	})

	t.Run("Case when the source of the Exchange is Layer API", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mockSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Background",
			Value: map[string]string{
				"from":   "UAH",
				"to":     "USD",
				"source": "Layer",
			},
		}
		var (
			to                      = mockSettingRes.Value["to"]
			from                    = mockSettingRes.Value["from"]
			source                  = mockSettingRes.Value["source"]
			yyyymmddDash            = currentTime.Format("2006-01-02")
			yyyymmddDashPreviousDay = yesterday.Format("2006-01-02")
		)

		fakeLayerResponse := `{
			"rates": {
				"USD": {
					"change": 0.88,
					"end_rate": 36.59
				}
			}
		}`

		mockStore.EXPECT().GetSettingByName(uid, models.Exchange).Return(*mockSettingRes, nil)

		httpmock.RegisterResponder(
			"GET",
			fmt.Sprintf(
				"https://api.apilayer.com/exchangerates_data/fluctuation?base=%s&start_date=%s&end_date=%s&symbols=%s",
				from,
				yyyymmddDashPreviousDay,
				yyyymmddDash,
				to,
			),
			httpmock.NewStringResponder(200, fakeLayerResponse),
		)

		exchange := service.GetExchange(id)
		assert.Equal(t, source, exchange.Source)
		assert.Equal(t, 0.88, exchange.Change)
		assert.Equal(t, 36.59, exchange.EndRate)
	})
}
