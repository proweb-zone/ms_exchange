package dto_properties

import "github.com/google/uuid"

type ProductPropertiesRequestDto struct {
	Data []PropertyDto `json:"data"`
}

type PropertyDto struct {
	Xml          uuid.UUID `json:"xml"`
	DeletionMark bool      `json:"deletion_mark"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
}

type ProductPropertyValuesRequestDto struct {
	Data []PropertyValuesDto `json:"data"`
}

type PropertyValuesDto struct {
	ProductXml uuid.UUID          `json:"product_xml"`
	Properties []DetailPropertDto `json:"properties"`
}

type DetailPropertDto struct {
	PropertyXml uuid.UUID `json:"xml_properties"`
	Value       string    `json:"value"`
}
