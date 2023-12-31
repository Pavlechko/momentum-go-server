package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Hashpassword string    `gorm:"type:varchar(255);not null"`
	Name         string    `gorm:"uniqueIndex;not null"`
	Settings     []Setting
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserInput struct {
	Name     string `json:"name" binding:"required,min=3,max=24"`
	Password string `json:"password" binding:"required,min=6,max=24"`
}

type UserResponseWithHash struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Hashpassword string    `json:"hashpassword"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}
