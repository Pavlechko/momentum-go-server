package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"momentum-go-server/internal/models"
	"momentum-go-server/internal/store"
	"momentum-go-server/internal/utils"
)

const minNameLeng = 3
const minPassLeng = 6

func CreateUser(user models.UserInput) (string, error) {
	err := validateUser(user)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	hashPasword, err := hashPassword(user.Password)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	user.Password = hashPasword

	newUser, err := store.CreateUser(user)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	token, tokenErr := utils.GenerateJWT(newUser)

	if tokenErr != nil {
		return "", fmt.Errorf("generation token error: %v", tokenErr.Error())
	}

	return token, nil
}

func GetUser(user models.UserInput) (string, error) {
	err := validateUser(user)

	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}

	existUser, err := store.GetUser(user)

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

func validateUser(user models.UserInput) error {
	if len(user.Name) < minNameLeng {
		return fmt.Errorf("your name is too short, min 3 symbols")
	} else if len(user.Password) < minPassLeng {
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

func VerifyPassword(hashedPassword, candidatePassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
}
