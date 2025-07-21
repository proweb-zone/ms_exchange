package dto_stocks

import "github.com/google/uuid"

type StocksRequestDto struct {
	General GeneralDto `json:"general"`
	Data    []DataDto  `json:"data"`
}

type GeneralDto struct {
	Storages []StorageDto `json:"storages"`
}

type StorageDto struct {
	Xml          uuid.UUID `json:"xml"`
	Name         string    `json:"name"`
	DeletionMark bool      `json:"deletion_mark"`
}

type DataDto struct {
	ProductXml uuid.UUID     `json:"product_xml"`
	Storages   []StoragesDto `json:"storages"`
}

type StoragesDto struct {
	StorageXml uuid.UUID `json:"storage_xml"`
	// DeletionMark bool      `json:"deletion_mark"`
	Quantity float64 `json:"quantity"`
}
