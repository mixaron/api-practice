package handler

import (
	"api-practice/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service}
}

func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	preview, err := c.FormFile("preview")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "preview file required")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["attachments"]

	userID := c.Locals("userID").(uint)

	article, err := h.service.CreateArticle(userID, title, content, preview, files)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   article,
	})
}
