package dto

type BaseResponse[T any] struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"operation done"`
	Data    *T     `json:"data,omitempty"`
}

// TokenBaseResponse, ErrorBaseResponse, SuccessRegisterBaseResponse для документации
type TokenBaseResponse struct {
	Status  string        `json:"status" example:"success"`
	Message string        `json:"message" example:"Authenticated"`
	Data    TokenResponse `json:"data"`
}

type ErrorBaseResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Error message"`
}

type SuccessRegisterBaseResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Success message"`
}

type SuccessProfileBaseResponse struct {
	Status  string          `json:"status" example:"success"`
	Message string          `json:"message" example:"Success message"`
	Data    ProfileResponse `json:"data"`
}
