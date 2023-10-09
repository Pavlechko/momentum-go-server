package store

import (
	"fmt"
	"log"
	"momentum-go-server/internal/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB
var err error

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
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("Connected successfully to the Database")

	DB.AutoMigrate(&models.User{})
}

// Creates a record of a new user in the database if his name is unique
func CreateUser(user models.UserInput) (*models.UserResponse, error) {
	var userModel models.User

	if isEntityExist(userModel, user.Name) {
		return &models.UserResponse{}, fmt.Errorf("user with name: %v already exists", user.Name)
	}

	now := time.Now()

	newUser := models.User{
		Name:         user.Name,
		Hashpassword: user.Password,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result := DB.Create(&newUser)

	if result.Error != nil {
		return &models.UserResponse{}, fmt.Errorf(result.Error.Error())
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	return userResponse, nil
}

func GetUser(user models.UserInput) (*models.UserResponse, error) {
	var userModel models.User

	if !isEntityExist(userModel, user.Name) {
		return &models.UserResponse{}, fmt.Errorf("user with the name: %v not found", user.Name)
	}

	result := DB.Find(&userModel, "name = ?", user.Name)

	if result.Error != nil {
		return &models.UserResponse{}, fmt.Errorf(result.Error.Error())
	}

	userResponse := &models.UserResponse{
		ID:        userModel.ID,
		Name:      userModel.Name,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
	return userResponse, nil

}

func isEntityExist(model models.User, name string) bool {

	entityExist := DB.Find(&model, "name = ?", name)

	return entityExist.RowsAffected == 1
}
