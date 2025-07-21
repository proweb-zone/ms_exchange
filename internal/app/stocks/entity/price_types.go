package entity

import (
	"time"

	"github.com/google/uuid"
)

type PriceTypes struct {
	Xml       uuid.UUID `json:"xml" gorm:"type:uuid;type:uuid;primary_key"`
	Name      string    `json:"name" gorm:"column:name"`
	IsActive  bool      `json:"is_active" gorm:"type:bool"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PriceTypes) TableName() string {
	return "price_types"
}
