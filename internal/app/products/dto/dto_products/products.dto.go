package dto_products

import "github.com/google/uuid"

type PreparedData[T any] struct {
	CreateList                []T
	UpdateList                []T
	UpdatedXmls               []uuid.UUID
	ForRepeatRequestCatB2cMap map[uuid.UUID]*T
	ForRepeatRequestCatB2bMap map[uuid.UUID]*T
	ForRepeatRequestBrandMap  map[uuid.UUID]*T
}
