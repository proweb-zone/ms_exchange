package dto_products

import "github.com/google/uuid"

type ProductsRequestDto struct {
	General GeneralDto   `json:"general"`
	Data    []ProductDto `json:"data"`
}

type GeneralDto struct {
	Units []UnitsDto `json:"units"`
}

type UnitsDto struct {
	Xml  string `json:"xml"`
	Name string `json:"name"`
}

type ProductDto struct {
	ProductXml   uuid.UUID  `json:"product_xml"`
	DeletionMark bool       `json:"deletion_mark"`
	Name         string     `json:"name"`
	Article      string     `json:"article"`
	CodeToledo   string     `json:"code_toledo"`
	Description  string     `json:"description"`
	BrandXml     string     `json:"brand"`
	UnitXml      string     `json:"unit"`
	Step         uint       `json:"step"`
	Categories   Categories `json:"categories"`
}

type Categories struct {
	CategoryXmlB2c uuid.UUID `json:"category_id_b2c"`
	CategoryXmlB2b uuid.UUID `json:"category_id_b2b"`
}
