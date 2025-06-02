package handler

import (
	. "api-practice/internal/dto"
	"api-practice/internal/service"
	"api-practice/wsocket"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	service service.ArticleService
	hub     wsocket.Server
}

func NewArticleHandler(service service.ArticleService, hub wsocket.Server) *ArticleHandler {
	return &ArticleHandler{service, hub}
}

// CreateArticle godoc
// @Summary Create Article post
// @Description Create Article by auth user. post http method
// @Tags article
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Param user body dto.ArticleRequest true "article create data"
// @Produce json
// @Success 200 {object} dto.ArticleBaseResponse "Article published successfully"
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 401 {object} dto.ErrorBaseResponse
// @Router /api/articles [post]
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

// AllArticles godoc
// @Summary Get all articles
// @Description find and return all users articles
// @Tags article
// @Produce json
// @Success 200 {object} dto.ArticleBaseResponse "Article published successfully"
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 401 {object} dto.ErrorBaseResponse
// @Router /api/articles [get]
func (h *ArticleHandler) AllArticles(c *fiber.Ctx) error {
	articles, err := h.service.GetAllArticles()
	if err != nil || len(articles) == 0 {
		return Error(c, "Published articles not found", 404)
	}
	res := ToArrayResponse(articles)

	return Success(c, "all articles", &res, fiber.StatusOK)
}

// PublishArticle godoc
// @Summary Publish article
// @Description user publish his own article
// @Tags article
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} dto.SuccessBaseResponse "Article published successfully"
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 401 {object} dto.ErrorBaseResponse
// @Router /api/articles/id [patch]
func (h *ArticleHandler) PublishArticle(c *fiber.Ctx) error {
	articleID := c.Params("id")
	userID := c.Locals("userID").(uint)
	err := h.service.PublishArticle(articleID, userID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}
	h.hub.SendMessage("New article: " + articleID)
	return SuccessNoData(c, "Article successfully published", fiber.StatusOK)
}

// DeleteArticle godoc
// @Summary Delete article
// @Description User deletes his own article by ID. Requires authentication.
// @Tags article
// @Security ApiKeyAuth
// @Param id path string true "Article ID" format(uuid)
// @Produce json
// @Success 200 {object} dto.SuccessBaseResponse "Article deleted successfully"
// @Failure 403 {object} dto.ErrorBaseResponse "Article does not belong to user"
// @Failure 404 {object} dto.ErrorBaseResponse "Article not found"
// @Router /api/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	articleID := c.Params("id")
	userID := c.Locals("userID").(uint)
	err := h.service.DeleteArticle(articleID, userID)
	if err != nil {
		return Error(c, err.Error(), fiber.StatusBadRequest)
	}
	return SuccessNoData(c, "Article successfully deleted", fiber.StatusOK)
}

// UpdateArticle godoc
// @Summary Update article
// @Description user update his own article
// @Tags article
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param user body dto.ArticleRequest true "article update data"
// @Success 200 {object} dto.ArticleBaseResponse "Article updated successfully"
// @Failure 400 {object} dto.ErrorBaseResponse
// @Failure 401 {object} dto.ErrorBaseResponse
// @Router /api/articles/id [put]
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

	return Success(c, "successfully updated", &res, fiber.StatusOK)
}
