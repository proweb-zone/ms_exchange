package entity

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type Products struct {
	Id             uint       `json:"id" gorm:"primaryKey"`
	Xml            uuid.UUID  `json:"xml" gorm:"column:xml;type:uuid;uniqueIndex"`
	CategoryXmlB2c *uuid.UUID `json:"category_xml_b2c" gorm:"column:category_xml_b2c;type:uuid;index"`
	CategoryXmlB2b *uuid.UUID `json:"category_xml_b2b" gorm:"column:category_xml_b2b;type:uuid;index"`
	IsActive       bool       `json:"is_active" gorm:"column:is_active;type:bool"`
	Name           string     `json:"name" gorm:"column:name"`
	Article        string     `json:"article" gorm:"column:article;size:50"`
	CodeToledo     string     `json:"code_toledo" gorm:"column:code_toledo"`
	Description    string     `json:"description" gorm:"column:description"`
	BrandXml       string     `json:"brand" gorm:"column:brand_xml"`
	Unit           string     `json:"unit" gorm:"column:unit"`
	Step           uint       `json:"step" gorm:"column:step"`

	CategoriesB2c ProductCategoriesB2c `json:"categories_b2c" gorm:"foreignKey:CategoryXmlB2c;references:Xml"`
	CategoriesB2b ProductCategoriesB2b `json:"categories_b2b" gorm:"foreignKey:CategoryXmlB2b;references:Xml"`
}

func (Products) TableName() string {
	return "products"
}
