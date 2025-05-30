package handler

import (
	. "api-practice/internal/dto"
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
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return Error(c, "user not authenticated", fiber.StatusUnauthorized)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return Error(c, "failed to parse form data", fiber.StatusBadRequest)
	}

	request := ArticleRequest{
		Title:      c.FormValue("title"),
		Content:    c.FormValue("content"),
		PreviewURL: "",
	}

	previewFile, err := c.FormFile("preview")
	if err != nil {
		return Error(c, "preview file required", fiber.StatusBadRequest)
	}

	attachmentFiles := form.File["attachments"]

	article, err := h.service.CreateArticle(userID, request.Title, request.Content, previewFile, attachmentFiles)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	res := ToResponse(article)

	return Success(c, "successfully created", &res, fiber.StatusCreated)
}

func (h *ArticleHandler) AllArticles(c *fiber.Ctx) error {
	articles, err := h.service.GetAllArticles()
	if err != nil || len(articles) == 0 {
		return Error(c, "Published articles not found", 404)
	}
	res := ToArrayResponse(articles)

	return Success(c, "all articles", &res, fiber.StatusOK)
}

func (h *ArticleHandler) PublishArticle(c *fiber.Ctx) error {
	articleID := c.Params("id")
	userID := c.Locals("userID").(uint)
	err := h.service.PublishArticle(articleID, userID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	return SuccessNoData(c, "Article successfully published", fiber.StatusOK)
}

func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	articleID := c.Params("id")
	userID := c.Locals("userID").(uint)
	err := h.service.DeleteArticle(articleID, userID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}

	return SuccessNoData(c, "Article successfully deleted", fiber.StatusOK)
}

func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return Error(c, "user not authenticated", fiber.StatusUnauthorized)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return Error(c, "failed to parse form data", fiber.StatusBadRequest)
	}

	request := ArticleRequest{
		Title:      c.FormValue("title"),
		Content:    c.FormValue("content"),
		PreviewURL: "",
	}

	previewFile, err := c.FormFile("preview")
	if err != nil {
		return Error(c, "preview file required", fiber.StatusBadRequest)
	}

	attachmentFiles := form.File["attachments"]
	articleID := c.Params("id")
	article, err := h.service.UpdateArticle(articleID, userID, request.Title, request.Content, previewFile, attachmentFiles)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusInternalServerError)
	}

	res := ToResponse(article)

	return Success(c, "successfully updated", &res, fiber.StatusCreated)
}
