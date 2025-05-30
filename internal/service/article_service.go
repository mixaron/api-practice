package service

import (
	"api-practice/internal/minio_service"
	"api-practice/internal/model"
	"api-practice/internal/repository"
	"mime/multipart"
)

type ArticleService interface {
	CreateArticle(userID uint, title, content string, preview *multipart.FileHeader, attachments []*multipart.FileHeader) (*model.Article, error)
}

type articleService struct {
	repo  repository.ArticleRepository
	minio minio_service.UploadService
}

func NewArticleService(repo repository.ArticleRepository, minio minio_service.UploadService) ArticleService {
	return &articleService{repo, minio}
}

func (s *articleService) CreateArticle(userID uint, title, content string, preview *multipart.FileHeader, attachments []*multipart.FileHeader) (*model.Article, error) {
	previewFile, err := preview.Open()
	if err != nil {
		return nil, err
	}
	defer previewFile.Close()

	previewURL, err := s.minio.UploadFile("articles", "preview/"+preview.Filename, previewFile, preview.Size, preview.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	article := &model.Article{
		Title:      title,
		Content:    content,
		PreviewURL: previewURL,
		UserID:     userID,
	}

	for _, file := range attachments {
		if len(attachments) > 10 {
			break
		}
		f, err := file.Open()
		if err != nil {
			continue
		}
		defer f.Close()

		url, err := s.minio.UploadFile("articles", "attachments/"+file.Filename, f, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			continue
		}

		article.Attachments = append(article.Attachments, model.Attachment{
			URL:      url,
			FileName: file.Filename,
		})
	}

	if err := s.repo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}
