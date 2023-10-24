package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"momentum-go-server/internal/models"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func GetPexelsBackgroundImage() models.FrontendBackgroundImageResponse {
	var response models.PexelsImageResponse

	client := &http.Client{}

	randomPage := rand.Intn(5000)

	apiURL := fmt.Sprintf("https://api.pexels.com/v1/search?query=Nature&page=%d&orientation=landscape&per_page=1", randomPage)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
	}

	req.Header.Add("Authorization", os.Getenv("PEXELS_BACKGROUND_IMAGE"))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading HTTP response body:", err)
	}

	json.Unmarshal([]byte(body), &response)

	frontendResponse := models.FrontendBackgroundImageResponse{
		Photographer: response.Photos[0].Photographer,
		Image:        response.Photos[0].Image.Original,
		Alt:          response.Photos[0].Alt,
		Sourse:       "pexels.com",
		SourseUrl:    response.Photos[0].SourseUrl,
	}
	return frontendResponse
}
