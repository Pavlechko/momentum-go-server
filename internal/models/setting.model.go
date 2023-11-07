package models

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	ID        int               `gorm:"type:int;primary_key"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null"`
	Name      string            `gorm:"type:varchar(255);not null"`
	Value     map[string]string `gorm:"type:jsonb;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
