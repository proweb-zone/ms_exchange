package entity

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type ProductPropertyTypes struct {
	Id       uint      `json:"id" gorm:"primaryKey"`
	Xml      uuid.UUID `json:"xml" gorm:"column:xml;type:uuid;uniqueIndex"`
	Name     string    `json:"name" gorm:"column:name"`
	Type     string    `json:"type" gorm:"column:type"`
	IsActive bool      `json:"is_active" gorm:"column:is_active;type:bool"`
}

func (ProductPropertyTypes) TableName() string {
	return "product_property_types"
}
