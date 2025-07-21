package entity

import (
	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type ProductPropertyValues struct {
	Id          uint      `json:"id" gorm:"primaryKey"`
	PropertyXml uuid.UUID `json:"property_xml" gorm:"column:property_xml;type:uuid;index"`
	ProductXml  uuid.UUID `json:"product_xml" gorm:"column:product_xml;type:uuid;index"`
	Value       string    `json:"value" gorm:"column:value"`

	PriceType ProductPropertyTypes `json:"property_types" gorm:"foreignKey:PropertyXml;references:Xml;constraint:OnDelete:CASCADE"`
	Product   Products             `json:"product" gorm:"foreignKey:ProductXml;references:Xml;constraint:OnDelete:CASCADE"`
}

func (ProductPropertyValues) TableName() string {
	return "product_property_values"
}
