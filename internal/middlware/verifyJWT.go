package middlware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {

			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

			if len(authHeader) != 2 {
				log.Println("20 line in verify JWT")
				handleError(w, "Malformed Token")
			} else {
				jwtToken := authHeader[1]
				parsedToken, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
					_, ok := t.Method.(*jwt.SigningMethodHMAC)

					if !ok {
						w.WriteHeader(http.StatusUnauthorized)
						_, err := w.Write([]byte("You're Unauthorized!"))
						if err != nil {
							return nil, err
						}
					}
					return []byte(os.Getenv("API_SECRET")), nil
				})

				if err != nil {
					handleError(w, "You're Unauthorized due to error parsing the JWT ")
				}

				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if ok {
					expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
					if !expirationTime.After(time.Now()) {
						handleError(w, "Token is expired")
					}
				} else {
					handleError(w, "Invalid token claims")
				}

				if parsedToken.Valid {
					endpointHandler(w, r)
				} else {
					handleError(w, "You're Unauthorized due to invalid token")
				}
			}

		} else {
			handleError(w, "You're Unauthorized due to No token in the header")
		}
	}
}

func handleError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte(message))
	if err != nil {
		return
	}
}
