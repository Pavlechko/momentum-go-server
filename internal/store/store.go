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
	settings := utils.GetDefaultSettings(id)

	var result = DB.Create(settings)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error creating default settings:", result.Error.Error())
		return
	}
}

func GetSettings(id uuid.UUID) models.SettingResponse {
	var settingsModel []models.Setting

	responseMap := make(map[string]models.ValueMap)

	result := DB.Find(&settingsModel, "user_id = ?", id)

	utils.InfoLogger.Println(result)
	utils.InfoLogger.Println(settingsModel)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user settings", result.Error.Error())
	}

	for _, setting := range settingsModel {
		responseMap[setting.Name] = setting.Value
	}

	settingRes := models.SettingResponse{
		Weather:    responseMap["Weather"],
		Quote:      responseMap["Quote"],
		Background: responseMap["Background"],
		Exchange:   responseMap["Exchange"],
		Market:     responseMap["Market"],
	}

	return settingRes
}

func UpdateSetting(id uuid.UUID, name models.SettingType, v map[string]string) (models.Setting, error) {
	var setting models.Setting

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		utils.ErrorLogger.Println("Error creating the transaction:", err)
		return setting, err
	}

	result := DB.Find(&setting, "user_id = ? AND name = ?", id, name)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user setting:", result.Error.Error())
		tx.Rollback()
		return setting, result.Error
	}

	result = DB.Model(&setting).Updates(models.Setting{Value: v})
	if result.Error != nil {
		utils.ErrorLogger.Println("Error update setting:", result.Error.Error())
		tx.Rollback()
		return setting, result.Error
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLogger.Println("Error transaction commit :", err)
		return setting, result.Error
	}

	return setting, nil
}
