package middleware

import (
	"api-practice/internal/auth"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(tokenService auth.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = c.Query("token")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
			}
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID in token",
			})
		}

		userID := uint(userIDFloat)
		c.Locals("userID", userID)

		return c.Next()
	}
}
