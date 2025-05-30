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
	"api-practice/internal/minio_service"
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
	dbErr := db.DB.AutoMigrate(&model.User{}, &model.Article{}, &model.Attachment{})
	if dbErr != nil {
		return
	}

	app := fiber.New()

	minioInit := minio_service.NewClientInitService()
	minioClient, err := minioInit.Init()
	if err != nil {
		log.Fatal(err)
	}

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository)
	tokenService := auth.NewTokenService(os.Getenv("SECRET"))
	userHandler := handler.NewUserHandler(userService, tokenService)
	profileHandler := handler.NewProfileHandler(userService)
	articleRepository := repository.NewArticleRepository(db.DB)
	minioService := minio_service.NewUploadService(minioClient)
	articleService := service.NewArticleService(articleRepository, minioService)
	articleHandler := handler.NewArticleHandler(articleService)

	routes.SetupRoutes(app, userHandler, profileHandler, tokenService, articleHandler)

	errApp := app.Listen(":3000")
	if errApp != nil {
		return
	}
}

// исправить названия по google style
// переделать указатели
// пересмотреть использование minio
