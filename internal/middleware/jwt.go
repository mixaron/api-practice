package middleware

import (
	"api-practice/internal/auth"
	"api-practice/internal/service"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func AuthMiddleware(tokenService auth.TokenService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			authHeader = c.Query("token")
			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).SendString("missing token")
			}
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok || !userService.IsUserExists(uint(userIDFloat)) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user ID in token",
			})
		}

		userID := uint(userIDFloat)
		c.Locals("userID", userID)

		return c.Next()
	}
}
