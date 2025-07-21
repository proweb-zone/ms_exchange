package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductPrices struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	IsActive   bool      `json:"is_active" gorm:"type:bool"`
	ProductXml uuid.UUID `json:"product_xml" gorm:"type:uuid;column:product_xml"`
	PriceXml   uuid.UUID `json:"price_xml" gorm:"type:uuid;column:price_xml"`
	Price      float64   `json:"price" gorm:"column:price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	PriceType PriceTypes `json:"type_price_xml" gorm:"foreignKey:PriceXml;references:Xml;constraint:OnDelete:CASCADE"`
}

func (ProductPrices) TableName() string {
	return "product_prices"
}
