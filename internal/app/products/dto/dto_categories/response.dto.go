package dto_categories

type UpdateBatchResponseDto struct {
	Updated     int64    `json:"updated"`
	Created     int64    `json:"created"`
	SoftDeleted int64    `json:"deleted"`
	Duplicated  []string `json:"duplicated"`
}
