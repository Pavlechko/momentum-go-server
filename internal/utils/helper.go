package utils

import (
	"log"
	"momentum-go-server/internal/models"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var InfoLogger *log.Logger
var ErrorLogger *log.Logger

func init() {
	path, err := filepath.Abs("./logs")
	if err != nil {
		log.Println("Error riding absolute path: ", err)
		return
	}

	myLog, err := os.OpenFile(path+"/momentum-log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
				"image":        "https://images.unsplash.com/photo-1465189684280-6a8fa9b19a7a?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w1MDU5Nzd8MHwxfHJhbmRvbXx8fHx8fHx8fDE2OTk2MjA3NDd8&ixlib=rb-4.0.3&q=80&w=1080",
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
