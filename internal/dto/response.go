package dto

import "github.com/gofiber/fiber/v2"

func Success[T any](c *fiber.Ctx, message string, data *T, status int) error {
	return c.Status(status).JSON(BaseResponse[T]{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessNoData(c *fiber.Ctx, message string, status int) error {
	return c.Status(status).JSON(BaseResponse[any]{
		Status:  "success",
		Message: message,
		Data:    nil,
	})
}

func Error(c *fiber.Ctx, message string, status int) error {
	return c.Status(status).JSON(BaseResponse[any]{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}
