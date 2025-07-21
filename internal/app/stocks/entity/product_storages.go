package entity

import (
	"time"

	"github.com/google/uuid"
)

type ProductStorages struct {
	Id         uint      `json:"id" gorm:"primaryKey"`
	ProductXml uuid.UUID `json:"product_xml" gorm:"type:uuid;column:product_xml"`
	StorageXml uuid.UUID `json:"storage_xml" gorm:"type:uuid;column:storage_xml"`
	Quantity   float64   `json:"quantity" gorm:"column:quantity"`
	IsActive   bool      `json:"is_active" gorm:"type:bool"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Storage Storages `json:"storage" gorm:"foreignKey:StorageXml;references:Xml;constraint:OnDelete:CASCADE"`
}

func (ProductStorages) TableName() string {
	return "product_storages"
}
