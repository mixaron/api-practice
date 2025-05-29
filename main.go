// @title Rest API для тестового задания
// @version 1.0
// @description API для авторизации по jwt, созданию статей, использования websocket
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"api-practice/internal/auth"
	"api-practice/internal/db"
	"api-practice/internal/handler"
	"api-practice/internal/model"
	"api-practice/internal/repository"
	"api-practice/internal/service"
	"api-practice/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	db.Init()
	dbErr := db.DB.AutoMigrate(&model.User{})
	if dbErr != nil {
		return
	}

	app := fiber.New()

	repo := repository.NewUserRepository()
	svc := service.NewUserService(repo)
	tokenService := auth.NewTokenService(os.Getenv("SECRET"))
	userHandler := handler.NewUserHandler(svc, tokenService)
	profileHandler := handler.NewProfileHandler(svc)
	routes.SetupRoutes(app, userHandler, profileHandler, tokenService)

	errApp := app.Listen(":3000")
	if errApp != nil {
		return
	}
}
