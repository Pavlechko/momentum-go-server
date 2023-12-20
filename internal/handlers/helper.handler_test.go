package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestIsDecodeJSONRequest(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockIService(ctrl)
	handler := &Handler{
		Service: mockService,
	}

	tests := []struct {
		caseName    string
		requestBody string
		expected    string
	}{
		{"Valid Body", `{"name": "username", "password": "password"}`, "username"},
		{"Invalid Body", `{"nam": "username", "password": "password"}`, ""},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			var user models.UserInput
			reqBody := strings.NewReader(test.requestBody)

			req := httptest.NewRequest("GET", "/", reqBody)
			w := httptest.NewRecorder()

			handler.IsDecodeJSONRequest(w, req, &user)

			assert.Equal(t, test.expected, user.Name)
		})
	}
}
