package dto_prices

import "github.com/google/uuid"

type PricesRequestDto struct {
	General GeneralDto `json:"general"`
	Data    []DataDto  `json:"data"`
}

type GeneralDto struct {
	PriceTypes []PriceTypeDto `json:"prices"`
}

type PriceTypeDto struct {
	Xml          uuid.UUID `json:"xml"`
	Name         string    `json:"name"`
	DeletionMark bool      `json:"deletion_mark"`
}

type DataDto struct {
	ProductXml uuid.UUID  `json:"product_xml"`
	Prices     []PriceDto `json:"prices"`
}

type PriceDto struct {
	TypePriceXml uuid.UUID `json:"type_price_xml"`
	Price        float64   `json:"price"`
}
