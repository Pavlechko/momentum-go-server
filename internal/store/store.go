package store

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

type Database struct {
	DB *gorm.DB
}

type Data interface {
	Close()
	CreateUser(user models.UserInput) (*models.UserResponse, error)
	GetUser(user models.UserInput) (*models.UserResponseWithHash, error)
	GetSettingByName(id uuid.UUID, name models.SettingType) (models.Setting, error)
	GetSettings(id uuid.UUID) models.SettingResponse
	UpdateSetting(id uuid.UUID, name models.SettingType, v map[string]string) (models.Setting, error)
}

var now = time.Now()

func ConnectDB() (*Database, error) {
	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Kyiv",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		utils.ErrorLogger.Println("Error reading HTTP response body:", err)
		return nil, err
	}
	fmt.Println("Connected successfully to the Database")
	utils.InfoLogger.Println("Connected successfully to the Database")

	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		utils.ErrorLogger.Println("failed to create 'uuid-ossp' extension")
	}
	fmt.Println("Extension 'uuid-ossp' created successfully.")
	utils.InfoLogger.Println("Extension 'uuid-ossp' created successfully.")

	err = db.AutoMigrate(&models.User{}, &models.Setting{})
	if err != nil {
		utils.ErrorLogger.Println("Error migration tables:", err)
		return nil, err
	}
	return &Database{DB: db}, err
}

func (db *Database) Close() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		utils.ErrorLogger.Println("Error closing database connection:", err)
		return
	}
	sqlDB.Close()
}

// Creates a record of a new user in the database if his name is unique
func (db *Database) CreateUser(user models.UserInput) (*models.UserResponse, error) {
	var userModel models.User

	if db.isEntityExist(&userModel, user.Name) {
		utils.ErrorLogger.Println("user with name:", user.Name, "already exists")
		return &models.UserResponse{}, fmt.Errorf("user with name: %v already exists", user.Name)
	}

	newUser := models.User{
		Name:         user.Name,
		Hashpassword: user.Password,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := db.DB.Create(&newUser)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error creating the user", user.Name, "Error message:", result.Error.Error())
		return &models.UserResponse{}, fmt.Errorf("%s", result.Error.Error())
	}

	userResponse := &models.UserResponse{
		ID:   newUser.ID,
		Name: newUser.Name,
	}
	db.createDefaultSettings(userResponse.ID)
	return userResponse, nil
}

func (db *Database) GetUser(user models.UserInput) (*models.UserResponseWithHash, error) {
	var userModel models.User

	if !db.isEntityExist(&userModel, user.Name) {
		utils.ErrorLogger.Println("user with the name:", user.Name, "not found")
		return &models.UserResponseWithHash{}, fmt.Errorf("user with the name: %v not found", user.Name)
	}

	result := db.DB.Find(&userModel, "name = ?", user.Name)

	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user", result.Error.Error())
		return &models.UserResponseWithHash{}, fmt.Errorf("%s", result.Error.Error())
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

func (db *Database) isEntityExist(model *models.User, name string) bool {
	entityExist := db.DB.Find(&model, "name = ?", name)

	return entityExist.RowsAffected == 1
}

// Creates a record with default settings for the user
func (db *Database) createDefaultSettings(id uuid.UUID) {
	settings := utils.GetDefaultSettings(id)

	var result = db.DB.Create(settings)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error creating default settings:", result.Error.Error())
		return
	}
}

func (db *Database) GetSettings(id uuid.UUID) models.SettingResponse {
	var settingsModel []models.Setting

	responseMap := make(map[string]models.ValueMap)

	result := db.DB.Find(&settingsModel, "user_id = ?", id)

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

func (db *Database) GetSettingByName(id uuid.UUID, name models.SettingType) (models.Setting, error) {
	var setting models.Setting

	// result := db.DB.Find(&setting, "user_id = ? AND name = ?", id, name)
	result := db.DB.Where("user_id = ? AND name = ?", id, name).Find(&setting)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user setting:", result.Error.Error())
		return setting, result.Error
	}
	return setting, nil
}

func (db *Database) UpdateSetting(id uuid.UUID, name models.SettingType, v map[string]string) (models.Setting, error) {
	var setting models.Setting

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		utils.ErrorLogger.Println("Error creating the transaction:", err)
		return setting, err
	}

	result := db.DB.Find(&setting, "user_id = ? AND name = ?", id, name)
	if result.Error != nil {
		utils.ErrorLogger.Println("Error finding user setting:", result.Error.Error())
		tx.Rollback()
		return setting, result.Error
	}

	result = db.DB.Model(&setting).Updates(models.Setting{Value: v})
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
