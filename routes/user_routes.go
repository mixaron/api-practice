package routes

import (
	_ "api-practice/docs"
	"api-practice/internal/auth"
	"api-practice/internal/handler"
	"api-practice/internal/middleware"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, h *handler.UserHandler, p *handler.ProfileHandler,
	ts auth.TokenService, a *handler.ArticleHandler) {

	api := app.Group("/api")
	api.Get("/swagger/*", fiberSwagger.WrapHandler)
	api.Get("/articles", a.AllArticles)

	authRoute := api.Group("/auth")
	authRoute.Post("/reg", h.Register)

	authRoute.Post("/verify", h.VerifyRegistration)

	authRoute.Post("/login", h.Authenticate)

	authGroup := api.Group("/", middleware.AuthMiddleware(ts))

	authGroup.Get("/profile", p.GetUserProfile)
	authGroup.Post("/articles", a.CreateArticle)
	authGroup.Patch("/articles/:id", a.PublishArticle)
	authGroup.Put("/articles/:id", a.UpdateArticle)
	authGroup.Delete("/articles/:id", a.DeleteArticle)

}
