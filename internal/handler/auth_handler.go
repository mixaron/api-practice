package handler

import (
	"api-practice/internal/auth"
	. "api-practice/internal/dto"
	"api-practice/internal/model"
	"api-practice/internal/service"
	"api-practice/internal/validator"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service      service.UserService
	tokenService auth.TokenService
}

func NewUserHandler(service service.UserService, tokenService auth.TokenService) *UserHandler {
	return &UserHandler{service: service, tokenService: tokenService}
}

// Register godoc
// @Summary Register a new user
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration info"
// @Success 201 {object} dto.SuccessRegisterBaseResponse
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 500 {object} dto.ErrorBaseResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return Error(c, "Invalid request body", fiber.StatusBadRequest)
	}

	if err := validator.Struct(req); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	user := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Register(&user); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	return SuccessNoData(c, "User successfully created", fiber.StatusCreated)
}

// Authenticate godoc
// @Summary Authenticate user
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.LoginRequest true "User login credentials"
// @Success 200 {object} dto.TokenBaseResponse
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 500 {object} dto.ErrorBaseResponse
// @Router /auth/login [post]
func (h *UserHandler) Authenticate(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	authUser, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	token, err := h.tokenService.GenerateToken(authUser.ID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	return Success(c, "Authenticated", &TokenResponse{Token: token}, fiber.StatusOK)
}
