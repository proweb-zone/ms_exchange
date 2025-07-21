package dto_products

type UpdateBatchResponseDto struct {
	Updated     int64 `json:"updated"`
	SoftDeleted int64 `json:"deleted"`
}
