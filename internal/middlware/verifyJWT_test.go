package middlware

import (
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVerifyJWT(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	os.Setenv("API_SECRET", "mySecretKey")

	uid, _ := uuid.NewRandom()

	user := &models.UserResponse{
		ID:   uid,
		Name: "TestUser",
	}

	t.Run("Valid Token", func(t *testing.T) {
		validToken, err := utils.GenerateJWT(user)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		w := httptest.NewRecorder()
		VerifyJWT(handler)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		validToken, err := utils.GenerateJWT(user)
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+validToken+"something")

		w := httptest.NewRecorder()
		VerifyJWT(handler)(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(
			t,
			"You're Unauthorized due to error parsing the JWT You're Unauthorized due to invalid token",
			w.Body.String(),
		)
	})

	t.Run("Malformed Token without Bearer teg", func(t *testing.T) {
		validToken, err := utils.GenerateJWT(user)
		assert.NoError(t, err)

		utils.InfoLogger.Println(validToken)
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", validToken)

		w := httptest.NewRecorder()
		VerifyJWT(handler)(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(
			t,
			"Malformed Token",
			w.Body.String(),
		)
	})

	t.Run("Expired Token", func(t *testing.T) {
		const expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDI4MjY1NTEsInVzZXJJZCI6IjIyZWU4NjYwLWM5Y2UtNGFjZS1iZDBmLTQ4ZjQ4NzhkNzZlOSIsInVzZXJOYW1lIjoiVGVzdFVzZXIifQ.jq0qABQxOZbo7UNl0obSWvHWuqhdOaP-J1LvzQYRkeE"

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		w := httptest.NewRecorder()
		VerifyJWT(handler)(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(
			t,
			"You're Unauthorized due to error parsing the JWT Token is expired",
			w.Body.String(),
		)
	})

	t.Run("No token in header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)

		w := httptest.NewRecorder()
		VerifyJWT(handler)(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(
			t,
			"You're Unauthorized due to No token in the header",
			w.Body.String(),
		)
	})
}
