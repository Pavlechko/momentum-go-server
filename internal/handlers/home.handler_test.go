package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
	"momentum-go-server/test/mocks"
)

func TestHome(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockService := mocks.NewMockIService(ctrl)
	handler := &Handler{
		Service: mockService,
	}

	resObj := &models.ResponseObj{}

	uid, _ := uuid.NewRandom()
	id := uid.String()

	user := &models.UserResponse{
		ID:   uid,
		Name: "TestUser",
	}

	fakeToken, err := utils.GenerateJWT(user)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fakeToken))

	w := httptest.NewRecorder()

	mockService.EXPECT().GetData(id).Return(*resObj).Times(1)

	handler.Home(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
