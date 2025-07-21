package dto_prices

type UpdateBatchResponseDto struct {
	Created     int64 `json:"created"`
	Updated     int64 `json:"updated"`
	SoftDeleted int64 `json:"deleted"`
}
