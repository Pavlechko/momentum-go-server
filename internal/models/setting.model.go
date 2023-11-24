package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SettingType string

const (
	Weather    SettingType = "Weather"
	Quote      SettingType = "Quote"
	Market     SettingType = "Market"
	Exchange   SettingType = "Exchange"
	Background SettingType = "Background"
)

type ValueMap map[string]string

type Setting struct {
	ID        int       `gorm:"type:int;primary_key"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_id_name,unique"`
	Name      string    `gorm:"type:varchar(255);not null;index:idx_user_id_name,unique"`
	Value     ValueMap  `gorm:"type:jsonb;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (v ValueMap) Value() (driver.Value, error) {
	return json.Marshal(v)
}

func (v *ValueMap) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, v)
	case string:
		return json.Unmarshal([]byte(src), v)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

type SettingResponse struct {
	Weather    ValueMap
	Quote      ValueMap
	Background ValueMap
	Exchange   ValueMap
	Market     ValueMap
}
