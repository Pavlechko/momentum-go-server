package store

import (
	"fmt"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB
var err error
var now = time.Now()

type SettingType string

const (
	Weather    SettingType = "Weather"
	Quote      SettingType = "Quote"
	Market     SettingType = "Market"
	Exchange   SettingType = "Exchange"
	Background SettingType = "Background"
)

func ContentDB() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Kyiv",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PORT"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return
	}
	fmt.Println("Connected successfully to the Database")
	utils.InfoLogger.Println("Connected successfully to the Database")

	err = DB.AutoMigrate(&models.User{}, &models.Setting{})
	if err != nil {
		utils.ErrorLogger.Println("Error migration tables:", err)
	}
}

// Creates a record of a new user in the database if his name is unique
func CreateUser(user models.UserInput) (*models.UserResponse, error) {
	var userModel models.User

	if isEntityExist(userModel, user.Name) {
		utils.ErrorLogger.Println("user with name:", user.Name, "already exists")
		return &models.UserResponse{}, fmt.Errorf("user with name: %v already exists", user.Name)
	}

	newUser := models.User{
		Name:         user.Name,
		Hashpassword: user.Password,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := DB.Create(&newUser)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error creating the user", user.Name, "Error message:", result.Error.Error())
		return &models.UserResponse{}, fmt.Errorf(result.Error.Error())
	}

	userResponse := &models.UserResponse{
		ID:   newUser.ID,
		Name: newUser.Name,
	}
	createDefaultSettings(userResponse.ID)
	return userResponse, nil
}

func GetUser(user models.UserInput) (*models.UserResponseWithHash, error) {
	var userModel models.User

	if !isEntityExist(userModel, user.Name) {
		utils.ErrorLogger.Println("user with the name:", user.Name, "not found")
		return &models.UserResponseWithHash{}, fmt.Errorf("user with the name: %v not found", user.Name)
	}

	result := DB.Find(&userModel, "name = ?", user.Name)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user", result.Error.Error())
		return &models.UserResponseWithHash{}, fmt.Errorf(result.Error.Error())
	}

	userResponse := &models.UserResponseWithHash{
		ID:           userModel.ID,
		Name:         userModel.Name,
		Hashpassword: userModel.Hashpassword,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}

	return userResponse, nil
}

func isEntityExist(model models.User, name string) bool {

	entityExist := DB.Find(&model, "name = ?", name)

	return entityExist.RowsAffected == 1
}

// Creates a record with default settings for the user
func createDefaultSettings(id uuid.UUID) {
	settings := []*models.Setting{
		{
			UserID: id,
			Name:   string(Weather),
			Value: map[string]string{
				"source": "OpenWeatherAPI",
				"city":   "Kyiv",
			},
		},
		{
			UserID: id,
			Name:   string(Background),
			Value: map[string]string{
				"source": "unsplash.com",
				"image":  "https://images.unsplash.com/photo-1472214103451-9374bd1c798e?crop=entropy&cs=tinysrgb&fit=max&fm=jpg&ixid=M3w1MDU5Nzd8MHwxfHJhbmRvbXx8fHx8fHx8fDE2OTkzNjgxODZ8&ixlib=rb-4.0.3&q=80&w=1080",
			},
		},
		{
			UserID: id,
			Name:   string(Quote),
			Value: map[string]string{
				"content": "The world makes way for the man who knows where he is going.",
				"author":  "Ralph Waldo Emerson",
			},
		},
		{
			UserID: id,
			Name:   string(Exchange),
			Value: map[string]string{
				"source": "NBU",
				"base":   "UAH",
			},
		},
		{
			UserID: id,
			Name:   string(Market),
			Value: map[string]string{
				"Symbol": "DAX",
			},
		},
	}

	var result = DB.Create(settings)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error creating default settings:", result.Error.Error())
		return
	}
}

func GetSettings(id uuid.UUID) {
	var settingsModel []models.Setting

	result := DB.Where(&settingsModel, "user_id = ?", id)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user settings", result.Error.Error())
	}
}

func UpdateSetting(id uuid.UUID, name SettingType, v map[string]string) {
	var setting models.Setting

	result := DB.Find(&setting, "user_id = ? AND name = ?", id, name)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user setting:", result.Error.Error())
		return
	}

	setting = models.Setting{
		Value: v,
	}
}
