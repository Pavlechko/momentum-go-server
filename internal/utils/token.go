package utils

import (
	"fmt"
	"momentum-go-server/internal/models"
	"net/http"
	"os"
	"strings"
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

func GetUserId(r *http.Request) string {
	if r.Header["Authorization"] != nil {

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeader) != 2 {
			ErrorLogger.Println("Malformed Token")
			return ""
		} else {
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
			} else {
				ErrorLogger.Println("Error parsing token")
				return ""
			}
		}
	} else {
		ErrorLogger.Println("You're Unauthorized due to No token in the header")
		return ""
	}
}
