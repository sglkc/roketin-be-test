package dto

type BaseResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// @name DataResponse
type DataResponse[T any] struct {
	BaseResponse
	Data T `json:"data"`
}

type PaginatedResponse[T any] struct {
	BaseResponse
	Data  []T `json:"data"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Count int `json:"count"`
}

type ErrorResponse struct {
	BaseResponse
	// Error string `json:"error"`
}
