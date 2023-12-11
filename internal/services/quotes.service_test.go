package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestGetRandomQuote(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	response := map[string]interface{}{
		"content": "Test quote",
		"author":  "Test author",
	}

	jsonResponse, _ := json.Marshal(response)

	httpmock.RegisterResponder(
		"GET",
		"https://api.quotable.io/random",
		httpmock.NewStringResponder(200, string(jsonResponse)),
	)

	res := GetRandomQuote()

	expected := models.QuoteResponse{
		Content: "Test quote",
		Author:  "Test author",
	}

	assert.Equal(t, expected, res)
}

func TestGetQuote(t *testing.T) {
	uid, _ := uuid.NewRandom()
	id := uid.String()
	currentTime := time.Now()
	previous12Hours := currentTime.Add(-12 * time.Hour)
	previous25Hours := currentTime.Add(-25 * time.Hour)

	mockStore := mocks.NewMockData(gomock.NewController(t))
	service := &Service{
		DB: mockStore,
	}

	t.Run("Case when the Quote was updated earlier than 24 hours", func(t *testing.T) {
		mockFreshSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Quote",
			Value: map[string]string{
				"content": "Some quote",
				"author":  "Some author",
			},
			UpdatedAt: previous12Hours,
		}

		mockStore.EXPECT().GetSettingByName(uid, models.Quote).Return(*mockFreshSettingRes, nil)
		quote := service.GetQuote(id)

		assert.Equal(t, "Some quote", quote.Content)
		assert.Equal(t, "Some author", quote.Author)
	})

	t.Run("Case when the Quote was updated more than 24 hours ago", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		response := map[string]string{
			"content": "New quote",
			"author":  "New author",
		}

		jsonResponse, _ := json.Marshal(response)

		httpmock.RegisterResponder(
			"GET",
			"https://api.quotable.io/random",
			httpmock.NewStringResponder(200, string(jsonResponse)),
		)

		mockOldSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Quote",
			Value: map[string]string{
				"content": "Some quote",
				"author":  "Some author",
			},
			UpdatedAt: previous25Hours,
		}

		mockNewSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Quote",
			Value: map[string]string{
				"content": "New quote",
				"author":  "New author",
			},
			UpdatedAt: currentTime,
		}

		mockStore.EXPECT().GetSettingByName(uid, models.Quote).Return(*mockOldSettingRes, nil)
		mockStore.EXPECT().UpdateSetting(uid, models.Quote, response).Return(*mockNewSettingRes, nil)
		quote := service.GetQuote(id)
		assert.NotEqual(t, "Some quote", quote.Content)
		assert.NotEqual(t, "Some author", quote.Author)
		assert.Equal(t, mockNewSettingRes.Value["content"], quote.Content)
		assert.Equal(t, mockNewSettingRes.Value["author"], quote.Author)
	})
}
