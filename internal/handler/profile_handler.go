package handler

import (
	. "api-practice/internal/dto"
	"api-practice/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	service service.UserService
}

func NewProfileHandler(service service.UserService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Получить профиль текущего пользователя
// @Tags profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.SuccessProfileBaseResponse
// @Failure 404 {object} dto.ErrorBaseResponse
// @Router /profile [get]
func (h *ProfileHandler) GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.service.GetProfile(userID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusNotFound)
	}

	return Success(c, "User profile", &ProfileResponse{Email: user.Email}, fiber.StatusOK)
}
