package routes

import (
	_ "api-practice/docs"
	"api-practice/internal/auth"
	"api-practice/internal/handler"
	"api-practice/internal/middleware"
	"api-practice/wsocket"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"log"
)

func SetupRoutes(app *fiber.App, h *handler.UserHandler, p *handler.ProfileHandler,
	ts auth.TokenService, a *handler.ArticleHandler, hub wsocket.Server) {

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

	authGroup.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	authGroup.Use("/ws", WebSocketContextMiddleware())

	authGroup.Get("/ws", websocket.New(func(c *websocket.Conn) {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			log.Println("Failed to get userID from context")
			c.Close()
			return
		}

		hub.HandleWS(c, fmt.Sprintf("%d", userID))
	}))
}

func WebSocketContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("ctx", c)
		return c.Next()
	}
}
