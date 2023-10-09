package store

import (
	"fmt"
	"log"
	"momentum-go-server/internal/models"
	"os"
	"strings"
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
func CreateUser(user models.UserInput) *models.UserResponse {
	var userModel models.User

	now := time.Now()

	newUser := models.User{
		Name:         user.Name,
		Hashpassword: user.Password,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	userExist := DB.Find(&userModel, "name = ?", newUser.Name)

	if userExist.RowsAffected == 1 {
		fmt.Println("User with that name already exists")
	} else {

		result := DB.Create(&newUser)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			fmt.Println("User with that name already exists")
		} else if result.Error != nil {
			fmt.Println("Error: " + result.Error.Error())
		}
	}

	userResponse := &models.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	return userResponse
}

func GetUser() {
	// var payload *models.SignUpInput
	var user models.User
	mockName := "Test1"
	result := DB.Find(&user, "name = ?", mockName)

	if result.Error != nil || result.RowsAffected == 0 {
		fmt.Println("Invalid name")
	}

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	fmt.Println(userResponse)

}
