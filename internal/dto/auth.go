package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"example@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"password213"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"example@gmail.com"`
	Password string `json:"password" validate:"required,min=8" example:"password213"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
