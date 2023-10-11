package middlware

import (
	"fmt"
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
				fmt.Println("Malformed token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
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

				fmt.Println("Header", parsedToken.Header, "Claims", parsedToken.Claims, "IsValid: ", parsedToken.Valid)

				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					_, err2 := w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
					if err2 != nil {
						return
					}
				}
				// TO-DO
				// check the time of life
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if ok {
					expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
					if expirationTime.After(time.Now()) {
						fmt.Println("Token is valid and not expired.")
					} else {
						w.WriteHeader(http.StatusUnauthorized)
						_, err := w.Write([]byte("Token is expired."))
						if err != nil {
							return
						}
					}
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					_, err := w.Write([]byte("Invalid token claims."))
					if err != nil {
						return
					}
				}

				if parsedToken.Valid {
					endpointHandler(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
					if err != nil {
						return
					}
				}
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	}
}
