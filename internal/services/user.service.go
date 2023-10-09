package services

import (
	"fmt"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.UserInput) (*models.UserResponse, error) {

	err := validateUser(user)

	if err != nil {
		return &models.UserResponse{}, fmt.Errorf(err.Error())
	}

	hashPasword, err := hashPassword(user.Password)

	if err != nil {
		return &models.UserResponse{}, fmt.Errorf(err.Error())
	}

	user.Password = hashPasword

	result, err := store.CreateUser(user)

	if err != nil {
		return &models.UserResponse{}, fmt.Errorf(err.Error())
	}

	return result, nil
}

func GetUser(user models.UserInput) (*models.UserResponse, error) {
	err := validateUser(user)

	if err != nil {
		return &models.UserResponse{}, fmt.Errorf(err.Error())
	}

	existUser, err := store.GetUser(user)

	if err != nil {
		return &models.UserResponse{}, fmt.Errorf(err.Error())
	}

	err = VerifyPassword(existUser.Hashpassword, user.Password)
	if err != nil {
		return &models.UserResponse{}, fmt.Errorf("incorrect password: %v", err.Error())
	}

	result := &models.UserResponse{
		ID:        existUser.ID,
		Name:      existUser.Name,
		CreatedAt: existUser.CreatedAt,
		UpdatedAt: existUser.UpdatedAt,
	}
	return result, nil

}

func validateUser(user models.UserInput) error {
	if len(user.Name) < 3 {
		return fmt.Errorf("your name is too short, min 3 symbols")
	} else if len(user.Password) < 6 {
		return fmt.Errorf("your password is too short, min 6 symbols")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
