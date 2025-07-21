package entity

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type ProductCategoriesB2c struct {
	Id       uint      `json:"id" gorm:"primaryKey"`
	ParentId *uint     `json:"parent_id" gorm:"column:parent_id;index"`
	Xml      uuid.UUID `json:"xml" gorm:"column:xml;type:uuid;uniqueIndex"`
	IsActive bool      `json:"is_active" gorm:"type:bool"`
	Name     string    `json:"name" gorm:"column:name"`

	Parent   *ProductCategoriesB2c  `gorm:"foreignKey:ParentId;references:Id"`
	Children []ProductCategoriesB2c `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE"`
}

func (ProductCategoriesB2c) TableName() string {
	return "product_categories_b2c"
}

type ProductCategoriesB2b struct {
	Id       uint      `json:"id" gorm:"primaryKey"`
	ParentId *uint     `json:"parent_id" gorm:"column:parent_id;index"`
	Xml      uuid.UUID `json:"xml" gorm:"column:xml;type:uuid;uniqueIndex"`
	IsActive bool      `json:"is_active" gorm:"type:bool"`
	Name     string    `json:"name" gorm:"column:name"`

	Parent   *ProductCategoriesB2b  `gorm:"foreignKey:ParentId;references:Id"`
	Children []ProductCategoriesB2b `gorm:"foreignKey:ParentId;constraint:OnDelete:CASCADE"`
}

func (ProductCategoriesB2b) TableName() string {
	return "product_categories_b2b"
}
