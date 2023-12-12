package services

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func getPexelsBackgroundImage() models.FrontendBackgroundImageResponse {
	const maxPages = 5000
	var response models.PexelsImageResponse
	var frontendResponse models.FrontendBackgroundImageResponse

	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}

	client := &http.Client{}

	randomPage, err := rand.Int(rand.Reader, big.NewInt(maxPages))
	if err != nil {
		utils.ErrorLogger.Printf("Error creating random number: %s", err.Error())
	}

	apiURL := fmt.Sprintf("https://api.pexels.com/v1/search?query=Nature&page=%d&orientation=landscape&per_page=1", randomPage)

	req, err := http.NewRequest("GET", apiURL, http.NoBody)
	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}

	req.Header.Add("Authorization", os.Getenv("PEXELS_BACKGROUND_IMAGE"))

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

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("Error json.Unmarshal response:", err)
	}

	frontendResponse = models.FrontendBackgroundImageResponse{
		Photographer: response.Photos[0].Photographer,
		Image:        response.Photos[0].Image.Original,
		Alt:          response.Photos[0].Alt,
		Source:       "pexels.com",
		SourceURL:    response.Photos[0].SourseURL,
	}
	return frontendResponse
}

func getUnsplashBackgroundImage() models.FrontendBackgroundImageResponse {
	var response models.UnsplashImageResponse
	var frontendResponse models.FrontendBackgroundImageResponse

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random?query=nature&orientation=landscape", http.NoBody)
	if err != nil {
		utils.ErrorLogger.Println("Error creating HTTP request:", err)
		return frontendResponse
	}

	req.Header.Add("Accept-Version", "v1")
	req.Header.Add("Authorization", "Client-ID "+os.Getenv("UNSPLASH_BACKGROUND_IMAGE"))

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

	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.ErrorLogger.Println("Error json.Unmarshal response:", err)
	}

	frontendResponse = models.FrontendBackgroundImageResponse{
		Photographer: response.Photographer.Name,
		Image:        response.Image.Regular,
		Alt:          response.Alt,
		Source:       "unsplash.com",
		SourceURL:    response.SourceURL.Image,
	}
	return frontendResponse
}

func GetRandomBackground(source string) models.FrontendBackgroundImageResponse {
	if source == "unsplash.com" {
		return getUnsplashBackgroundImage()
	}
	return getPexelsBackgroundImage()
}

func (s *Service) GetBackgroundData(userID string) models.FrontendBackgroundImageResponse {
	const oneDayHours = 24
	var response models.FrontendBackgroundImageResponse
	currentTime := time.Now()
	id, _ := uuid.Parse(userID)

	res, err := s.DB.GetSettingByName(id, models.Background)
	if err != nil {
		utils.ErrorLogger.Println("Error finding Background setting:", err)
		return response
	}
	dif := currentTime.Sub(res.UpdatedAt).Hours()

	if dif < oneDayHours {
		response = models.FrontendBackgroundImageResponse{
			Photographer: res.Value["photographer"],
			Image:        res.Value["image"],
			Alt:          res.Value["alt"],
			Source:       res.Value["source"],
			SourceURL:    res.Value["source_url"],
		}
		return response
	}
	response = GetRandomBackground(res.Value["source"])
	return response
}
