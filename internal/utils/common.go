package utils

import (
	"net/http"

	"momentum-go-server/internal/models"

	"github.com/google/uuid"
)

type IUtils interface {
	GenerateJWT(user *models.UserResponse) (string, error)
	GetUserID(r *http.Request) string
	GetDefaultSettings(id uuid.UUID) []*models.Setting
}
