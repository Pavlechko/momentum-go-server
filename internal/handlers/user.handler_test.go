package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		caseName string
		name     string
		password string
		expected error
	}{
		{"Short Name", "Sh", "passWord", fmt.Errorf("your name is too short, min 3 symbols")},
		{"Short Password", "Name", "pass", fmt.Errorf("your password is too short, min 6 symbols")},
		{"Valid User Input", "Name", "passWord", nil},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			user := models.UserInput{
				Name:     test.name,
				Password: test.password,
			}

			err := validateUser(user)

			assert.Equal(t, test.expected, err)
		})
	}
}

func TestSignInHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockIService(ctrl)
	handler := &Handler{
		Service: mockService,
	}

	tests := []struct {
		caseName       string
		requestBody    string
		expectedStatus int
		token          string
		err            error
	}{
		{"ValidUser", `{"name": "username", "password": "password"}`, http.StatusOK, "your_expected_token", nil},
		{"InvalidUser", `{"name": "sh", "password": "password"}`, http.StatusConflict, "", nil},
		{"UnauthorizedUser",
			`{"name": "username", "password": "wrongpassword"}`,
			http.StatusUnauthorized,
			"",
			fmt.Errorf("incorrect password or name"),
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			reqBody := strings.NewReader(test.requestBody)
			req := httptest.NewRequest("POST", "/auth/signin", reqBody)
			req.Header.Set("Content-Type", "application/json")

			var userInput models.UserInput
			err := json.Unmarshal([]byte(test.requestBody), &userInput)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			mockService.EXPECT().GetUser(userInput).Return(test.token, test.err)

			w := httptest.NewRecorder()

			handler.SignInHandler(w, req)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}

func TestSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockIService(ctrl)
	handler := &Handler{
		Service: mockService,
	}

	tests := []struct {
		caseName       string
		requestBody    string
		expectedStatus int
		token          string
		err            error
	}{
		{"ValidUser", `{"name": "username", "password": "password"}`, http.StatusOK, "your_expected_token", nil},
		{"InvalidUser", `{"name": "sh", "password": "password"}`, http.StatusConflict, "", nil},
		{"UnauthorizedUser",
			`{"name": "username", "password": "wrongpassword"}`,
			http.StatusConflict,
			"",
			fmt.Errorf("incorrect password or name"),
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			reqBody := strings.NewReader(test.requestBody)
			req := httptest.NewRequest("POST", "/auth/signin", reqBody)
			req.Header.Set("Content-Type", "application/json")

			var userInput models.UserInput
			err := json.Unmarshal([]byte(test.requestBody), &userInput)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			mockService.EXPECT().CreateUser(userInput).Return(test.token, test.err)

			w := httptest.NewRecorder()

			handler.SignUpHandler(w, req)
			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
