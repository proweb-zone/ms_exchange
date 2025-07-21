package dto_categories

type ProductCategoriesRequestDto struct {
	General GeneralDto    `json:"general"`
	Data    []CategoryDto `json:"data"`
}

type GeneralDto struct {
	Clean bool `json:"Clean"`
}

type CategoryDto struct {
	Xml          string        `json:"xml"`
	DeletionMark bool          `json:"deletion_mark"`
	Name         string        `json:"name"`
	Parent       string        `json:"parent"`
	Children     []CategoryDto `json:"children"`
}
