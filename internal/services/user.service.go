package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/utils"
)

func (s *Service) CreateUser(user models.UserInput) (string, error) {
	hashPasword, err := hashPassword(user.Password)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	user.Password = hashPasword

	newUser, err := s.DB.CreateUser(user)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	token, tokenErr := utils.GenerateJWT(newUser)

	if tokenErr != nil {
		return "", fmt.Errorf("generation token error: %v", tokenErr.Error())
	}

	return token, nil
}

func (s *Service) GetUser(user models.UserInput) (string, error) {
	existUser, err := s.DB.GetUser(user)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	err = VerifyPassword(existUser.Hashpassword, user.Password)
	if err != nil {
		return "", fmt.Errorf("incorrect password or name. Error:%s", err.Error())
	}

	userData := &models.UserResponse{
		ID:   existUser.ID,
		Name: existUser.Name,
	}

	token, tokenErr := utils.GenerateJWT(userData)

	if tokenErr != nil {
		return "", fmt.Errorf("generation token error: %v", tokenErr.Error())
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
