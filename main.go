// @title Rest API для тестового задания
// @version 1.0
// @description API для авторизации по jwt, созданию статей, использования wsocket
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
	"api-practice/internal/minio_service"
	"api-practice/internal/repository"
	"api-practice/internal/service"
	"api-practice/routes"
	"api-practice/wsocket"
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

	dbErr := db.Init()
	if dbErr != nil {
		return
	}

	app := fiber.New()

	minioInit := minio_service.NewClientInitService()
	minioClient, err := minioInit.Init()
	if err != nil {
		log.Fatal(err)
	}

	websocket := wsocket.NewServer()
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	tokenService := auth.NewTokenService(os.Getenv("SECRET"))
	userHandler := handler.NewUserHandler(userService, tokenService)
	profileHandler := handler.NewProfileHandler(userService)
	articleRepository := repository.NewArticleRepository(db.DB)
	minioService := minio_service.NewUploadService(minioClient)
	articleService := service.NewArticleService(articleRepository, minioService)
	articleHandler := handler.NewArticleHandler(articleService, *websocket)
	routes.SetupRoutes(app, userHandler, profileHandler, tokenService, articleHandler, *websocket)

	errApp := app.Listen(":3000")
	if errApp != nil {
		return
	}
}

// todo исправить названия по google style
// todo переделать указатели
// todo пересмотреть использование minio
// todo Отдавать кастомные ошибки
// todo доделать профиль
