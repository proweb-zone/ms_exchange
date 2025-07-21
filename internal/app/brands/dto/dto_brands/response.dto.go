package dto_brands

type UpdateBatchResponseDto struct {
	Updated     int64    `json:"updated"`
	Created     int64    `json:"created"`
	Duplicated  []string `json:"duplicated"`
	SoftDeleted int64    `json:"deleted"`
}
