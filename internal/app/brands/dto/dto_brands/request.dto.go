package dto_brands

type BrandsRequestDto struct {
	Data []BrandDto `json:"data"`
}

type BrandDto struct {
	Xml          string `json:"xml"`
	DeletionMark bool   `json:"deletion_mark"`
	Name         string `json:"name"`
}
