package routes

import (
	_ "api-practice/docs"
	"api-practice/internal/auth"
	"api-practice/internal/handler"
	"api-practice/internal/middleware"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, h *handler.UserHandler, p *handler.ProfileHandler, ts auth.TokenService) {
	api := app.Group("/api")
	api.Get("/swagger/*", fiberSwagger.WrapHandler)

	authRoute := api.Group("/auth")
	authRoute.Post("/register", h.Register)
	authRoute.Post("/login", h.Authenticate)

	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware(ts))
	profile.Get("", p.GetUserProfile)
}
