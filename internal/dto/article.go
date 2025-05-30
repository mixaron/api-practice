package dto

import (
	"api-practice/internal/model"
	"github.com/samber/lo"
	"time"
)

type ArticleRequest struct {
	Title       string `gorm:"size:255"`
	Content     string `gorm:"type:text"`
	PreviewURL  string
	Attachments []AttachmentRequest
}
type ArticleResponse struct {
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	PreviewURL  string               `json:"preview_url"`
	isPublished bool                 `json:"is_published"`
	AuthorID    uint                 `json:"author_id"`
	CreatedAt   time.Time            `json:"created_at"`
	Attachments []AttachmentResponse `json:"attachments"`
}

func ToResponse(article *model.Article) ArticleResponse {
	return ArticleResponse{
		Title:       article.Title,
		Content:     article.Content,
		PreviewURL:  article.PreviewURL,
		isPublished: article.IsPublished,
		AuthorID:    article.User.ID,
		CreatedAt:   article.CreatedAt,
		Attachments: lo.Map(article.Attachments, func(a model.Attachment, _ int) AttachmentResponse {
			return AttachmentResponse{
				URL: a.URL,
			}
		}),
	}
}

func ToArrayResponse(articles []model.Article) []ArticleResponse {
	if articles == nil {
		return []ArticleResponse{}
	}

	res := make([]ArticleResponse, len(articles))
	for i, article := range articles {
		res[i] = ArticleResponse{
			Title:       article.Title,
			Content:     article.Content,
			PreviewURL:  article.PreviewURL,
			isPublished: article.IsPublished,
			AuthorID:    article.User.ID,
			CreatedAt:   article.CreatedAt,
			Attachments: lo.Map(article.Attachments, func(a model.Attachment, _ int) AttachmentResponse {
				return AttachmentResponse{
					URL: a.URL,
				}
			}),
		}
	}
	return res
}
