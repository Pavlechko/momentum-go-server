package services

import (
	"fmt"
	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.UserInput) (string, error) {

	err := validateUser(user)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	hashPasword, err := hashPassword(user.Password)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	user.Password = hashPasword

	newUser, err := store.CreateUser(user)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	token, tokenErr := utils.GenerateJWT(newUser)

	if tokenErr != nil {
		return "", fmt.Errorf("generation token error: %v", err.Error())
	}

	return token, nil
}

func GetUser(user models.UserInput) (string, error) {
	err := validateUser(user)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	existUser, err := store.GetUser(user)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	err = VerifyPassword(existUser.Hashpassword, user.Password)
	if err != nil {
		return "", fmt.Errorf("incorrect password or name")
	}

	userData := &models.UserResponse{
		ID:   existUser.ID,
		Name: existUser.Name,
	}

	token, tokenErr := utils.GenerateJWT(userData)

	if tokenErr != nil {
		return "", fmt.Errorf("generation token error: %v", err.Error())
	}

	return token, nil
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
