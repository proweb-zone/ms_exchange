package dto

type Result[T any] struct {
	Result T `json:"result"`
}

type ErrorResponse struct {
	Result any    `json:"result"`
	Error  string `json:"error"`
}
