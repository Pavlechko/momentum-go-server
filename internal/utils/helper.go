package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"momentum-go-server/internal/models"
)

const chmod = 0666

var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func init() {
	pathSeck := filepath.FromSlash("logs/momentum-log.log")
	myLog, err := os.OpenFile(pathSeck, os.O_RDWR|os.O_CREATE|os.O_APPEND, chmod)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}
	InfoLogger = log.New(myLog, "[Info]:\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	ErrorLogger = log.New(myLog, "[Error]:\t", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}

func GetDefaultSettings(id uuid.UUID) []*models.Setting {
	defaultSettings := []*models.Setting{
		{
			UserID: id,
			Name:   string(models.Weather),
			Value: map[string]string{
				"source": "OpenWeather",
				"city":   "Kyiv",
			},
		},
		{
			UserID: id,
			Name:   string(models.Background),
			Value: map[string]string{
				"source":       "unsplash.com",
				"image":        "https://images.unsplash.com/photo-1465189684280-6a8fa9b19a7a?q=80&w=1080",
				"photographer": "Kalen Emsley",
				"alt":          "body of water surrounding with trees",
				"source_url":   "https://unsplash.com/photos/body-of-water-surrounding-with-trees-_LuLiJc1cdo",
			},
		},
		{
			UserID: id,
			Name:   string(models.Quote),
			Value: map[string]string{
				"content": "The world makes way for the man who knows where he is going.",
				"author":  "Ralph Waldo Emerson",
			},
		},
		{
			UserID: id,
			Name:   string(models.Exchange),
			Value: map[string]string{
				"source": "NBU",
				"from":   "UAH",
				"to":     "USD",
			},
		},
		{
			UserID: id,
			Name:   string(models.Market),
			Value: map[string]string{
				"symbol": "DAX",
			},
		},
	}
	return defaultSettings
}
