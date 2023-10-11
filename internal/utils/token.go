package utils

import (
	"fmt"
	"momentum-go-server/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

func GenerateJWT(user *models.UserResponse) (string, error) {
	now := time.Now().UTC()
	privateKey := os.Getenv("API_SECRET")

	// timeExpires := os.Getenv("TOKEN_EXPIRED_IN")
	claims := &jwt.MapClaims{
		"exp":      now.Add(time.Duration(60) * time.Minute).Unix(),
		"userId":   user.ID,
		"userName": user.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(privateKey))

	if err != nil {
		return "", fmt.Errorf("generation token error: %v", err.Error())
	}

	return tokenString, err
}
