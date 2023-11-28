package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"momentum-go-server/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const sixty int8 = 60
const correctHeaderLength = 2

func GenerateJWT(user *models.UserResponse) (string, error) {
	err := godotenv.Load()
	if err != nil {
		ErrorLogger.Printf("Error loading .env file, %s", err.Error())
	}
	now := time.Now().UTC()
	privateKey := os.Getenv("API_SECRET")

	claims := &jwt.MapClaims{
		"exp":      now.Add(time.Duration(sixty) * time.Minute).Unix(),
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

func GetUserID(r *http.Request) string {
	if r.Header["Authorization"] != nil {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != correctHeaderLength {
			ErrorLogger.Println("Malformed Token")
			return ""
		}

		jwtToken := authHeader[1]
		parsedToken, _ := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("API_SECRET")), nil
		})

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if ok {
			id, ok := claims["userId"].(string)
			if !ok {
				ErrorLogger.Println("Error parsing userId")
			}
			return id
		}
		ErrorLogger.Println("Error parsing token")
		return ""
	}
	ErrorLogger.Println("You're Unauthorized due to No token in the header")
	return ""
}
