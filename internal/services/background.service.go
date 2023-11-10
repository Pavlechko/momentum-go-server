package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetBackgroundData() models.BackgroundData {

	pexelsRes := getPexelsBackgroundImage()
	unsplashRes := getUnsplashBackgroundImage()

	Backgrounds := models.BackgroundData{
		Pexels:   pexelsRes,
		Unsplash: unsplashRes,
	}

	return Backgrounds
}

func getPexelsBackgroundImage() models.FrontendBackgroundImageResponse {
	var response models.PexelsImageResponse
	var frontendResponse models.FrontendBackgroundImageResponse

	client := &http.Client{}

	randomPage := rand.Intn(5000)

	apiURL := fmt.Sprintf("https://api.pexels.com/v1/search?query=Nature&page=%d&orientation=landscape&per_page=1", randomPage)

	req, err := http.NewRequest("GET", apiURL, nil)
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

	json.Unmarshal([]byte(body), &response)

	frontendResponse = models.FrontendBackgroundImageResponse{
		Photographer: response.Photos[0].Photographer,
		Image:        response.Photos[0].Image.Original,
		Alt:          response.Photos[0].Alt,
		Source:       "pexels.com",
		SourceUrl:    response.Photos[0].SourseUrl,
	}
	return frontendResponse
}

func getUnsplashBackgroundImage() models.FrontendBackgroundImageResponse {
	var response models.UnsplashImageResponse
	var frontendResponse models.FrontendBackgroundImageResponse

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.unsplash.com/photos/random?query=nature&orientation=landscape", nil)
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

	json.Unmarshal([]byte(body), &response)

	frontendResponse = models.FrontendBackgroundImageResponse{
		Photographer: response.Photographer.Name,
		Image:        response.Image.Regular,
		Alt:          response.Alt,
		Source:       "unsplash.com",
		SourceUrl:    response.SourceUrl.Image,
	}
	return frontendResponse
}

func GetRandomBackground() models.FrontendBackgroundImageResponse {
	randomInt := rand.Intn(2)
	if randomInt == 1 {
		return getUnsplashBackgroundImage()
	}
	return getPexelsBackgroundImage()
}
