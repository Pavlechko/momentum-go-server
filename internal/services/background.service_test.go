package services

import (
	"encoding/json"
	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetBackgroundData(t *testing.T) {
	uid, _ := uuid.NewRandom()
	id := uid.String()
	currentTime := time.Now()
	previous12Hours := currentTime.Add(-12 * time.Hour)
	previous25Hours := currentTime.Add(-25 * time.Hour)

	mockStore := mocks.NewMockData(gomock.NewController(t))
	service := &Service{
		DB: mockStore,
	}

	t.Run("Case when the Background was updated earlier than 24 hours", func(t *testing.T) {
		mockFreshSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Background",
			Value: map[string]string{
				"alt":          "Image",
				"image":        "https://images.unsplash.com/fresh",
				"source":       "unsplash.com",
				"source_url":   "https://unsplash.com/photos/fresh",
				"photographer": "Andreas Gücklhorn",
			},
			UpdatedAt: previous12Hours,
		}

		mockStore.EXPECT().GetSettingByName(uid, models.Background).Return(*mockFreshSettingRes, nil)
		background := service.GetBackgroundData(id)

		assert.Equal(t, mockFreshSettingRes.Value["image"], background.Image)
	})

	t.Run("Case when the Background was updated more than 24 hours ago", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		response := map[string]string{
			"alt":          "New Image",
			"image":        "https://images.unsplash.com/new",
			"source":       "unsplash.com",
			"source_url":   "https://unsplash.com/photos/new",
			"photographer": "New Photographer",
		}

		jsonResponse, _ := json.Marshal(response)

		httpmock.RegisterResponder(
			"GET",
			"https://api.unsplash.com/photos/random?query=nature&orientation=landscape",
			httpmock.NewStringResponder(200, string(jsonResponse)),
		)

		mockOldSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Background",
			Value: map[string]string{
				"alt":          "Image",
				"image":        "https://images.unsplash.com/Old",
				"source":       "unsplash.com",
				"source_url":   "https://unsplash.com/photos/Old",
				"photographer": "Andreas Gücklhorn",
			},
			UpdatedAt: previous25Hours,
		}

		mockNewSettingRes := &models.Setting{
			UserID: uid,
			Name:   "Background",
			Value: map[string]string{
				"alt":          "Image",
				"image":        "https://images.unsplash.com/New",
				"source":       "unsplash.com",
				"source_url":   "https://unsplash.com/photos/New",
				"photographer": "Andreas Gücklhorn",
			},
			UpdatedAt: currentTime,
		}

		mockStore.EXPECT().GetSettingByName(uid, models.Background).Return(*mockOldSettingRes, nil)
		mockStore.EXPECT().UpdateSetting(uid, models.Background, gomock.Any()).Return(*mockNewSettingRes, nil)
		background := service.GetBackgroundData(id)
		assert.NotEqual(t, mockOldSettingRes.Value["image"], background.Image)
		assert.Equal(t, mockNewSettingRes.Value["image"], background.Image)
	})
}

func TestGetRandomBackground(t *testing.T) {
	mockUnsplashRes := &models.FrontendBackgroundImageResponse{
		Image:  "https://images.unsplash.com/New",
		Source: "unsplash.com",
	}

	mockPixelsRes := &models.FrontendBackgroundImageResponse{
		Image:  "https://images.pexels.com/photos/New",
		Source: "pexels.com",
	}

	fakeUnsplashResponse := `{
		"user": {"name": "John Doe"},
		"urls": {"regular": "https://images.unsplash.com/New"},
		"alt_description": "Nature",
		"links": {"html": "https://unsplash.com/photos/New"}}`

	fakePixelsResponse := `{
		"photos": [
			{
				"photographer": "Darius",
				"src": {"landscape": "https://images.pexels.com/photos/New"},
				"alt": "Pixels Nature",
				"url": "https://unsplash.com/photos/New"
			}
		]}`

	t.Run("Case when the source of the Background is unsplash.com", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET",
			"https://api.unsplash.com/photos/random?query=nature&orientation=landscape",
			httpmock.NewStringResponder(200, fakeUnsplashResponse),
		)

		background := GetRandomBackground("unsplash.com")

		assert.Equal(t, mockUnsplashRes.Image, background.Image)
	})

	t.Run("Case when the source of the Background is pexels.com", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder(
			"GET",
			`=~https://api\.pexels\.com/v1/search\?.*?&orientation=landscape&per_page=1`,
			httpmock.NewStringResponder(200, fakePixelsResponse),
		)

		background := GetRandomBackground("pexels.com")

		assert.Equal(t, mockPixelsRes.Image, background.Image)
	})
}
