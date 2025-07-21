package entity

import "github.com/google/uuid"

type Brands struct {
	Id              uint      `json:"id" gorm:"primaryKey"`
	Xml             uuid.UUID `json:"xml" gorm:"column:xml;type:uuid;uniqueIndex"`
	IsActive        bool      `json:"is_active" gorm:"type:bool"`
	Name            string    `json:"name" gorm:"column:name"`
	Description     string    `json:"description" gorm:"column:description"`
	Preview         string    `json:"preview" gorm:"column:preview"`
	PreviewText     string    `json:"preview_text" gorm:"column:preview_text"`
	PreviewImage    string    `json:"preview_image" gorm:"column:preview_image"`
	PreviewImageBin []byte    `json:"preview_image_bin"`
	Sort            int       `json:"sort" gorm:"column:sort"`
	Images          []string  `json:"images" gorm:"serializer:json"`
	Tags            []string  `json:"tags" gorm:"serializer:json"`
}

func (Brands) TableName() string {
	return "brands"
}
