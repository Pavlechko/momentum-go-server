package services

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"momentum-go-server/internal/models"
	"momentum-go-server/test/mocks"
)

func TestGetUser(t *testing.T) {
	uid, _ := uuid.NewRandom()

	mockStore := mocks.NewMockData(gomock.NewController(t))

	service := &Service{
		DB: mockStore,
	}

	testUserInput := models.UserInput{
		Name:     "MockUser",
		Password: "Serhii",
	}

	mockUserRes := &models.UserResponseWithHash{
		ID:           uid,
		Name:         "MockUser",
		Hashpassword: "$2a$10$h4HREfn9FrYABbaHbMjv8OZedAfbUGBcJEHjn.d3qVp5d5IeRUbgW",
	}

	mockStore.EXPECT().GetUser(gomock.Eq(testUserInput)).Return(mockUserRes, nil)

	token, err := service.GetUser(testUserInput)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	userId, _ := claims["userId"].(string)
	id, _ := uuid.Parse(userId)
	userName, _ := claims["userName"].(string)

	assert.NoError(t, err)
	assert.Equal(t, id, mockUserRes.ID)
	assert.Equal(t, userName, mockUserRes.Name)
}

func TestCreateUser(t *testing.T) {
	uid, _ := uuid.NewRandom()
	mockStore := mocks.NewMockData(gomock.NewController(t))
	service := &Service{
		DB: mockStore,
	}

	testUserInput := models.UserInput{
		Name:     "NewUser",
		Password: "Serhii",
	}

	mockUserRes := &models.UserResponse{
		ID:   uid,
		Name: "NewUser",
	}

	mockStore.EXPECT().CreateUser(gomock.Any()).Return(mockUserRes, nil)

	token, err := service.CreateUser(testUserInput)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	parsedToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("API_SECRET")), nil
	})

	claims, _ := parsedToken.Claims.(jwt.MapClaims)
	userName, _ := claims["userName"].(string)

	assert.NoError(t, err)
	assert.Equal(t, userName, mockUserRes.Name)
}
