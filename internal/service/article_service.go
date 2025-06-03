package service

import (
	"api-practice/internal/minio_service"
	"api-practice/internal/model"
	"api-practice/internal/repository"
	"fmt"
	"mime/multipart"
	"time"
)

type ArticleService interface {
	CreateArticle(userID uint, title, content string, preview *multipart.FileHeader, attachments []*multipart.FileHeader) (*model.Article, error)
	GetAllArticles() ([]model.Article, error)
	PublishArticle(articleID string, userID uint) error
	DeleteArticle(articleID string, userID uint) error
	UpdateArticle(articleID string, userID uint, title, content string, preview *multipart.FileHeader,
		attachments []*multipart.FileHeader) (*model.Article, error)
	GetAllArticlesAfterTime(time time.Time) ([]model.Article, error)
}

type articleService struct {
	repo  repository.ArticleRepository
	minio *minio_service.UploadServiceImpl
}

func NewArticleService(repo repository.ArticleRepository, minio *minio_service.UploadServiceImpl) ArticleService {
	return &articleService{repo, minio}
}

func (s *articleService) CreateArticle(userID uint, title, content string, preview *multipart.FileHeader, attachments []*multipart.FileHeader) (*model.Article, error) {
	previewFile, err := preview.Open()
	if err != nil {
		return nil, err
	}
	defer previewFile.Close()

	previewURL, err := s.minio.UploadFile("articles1", "preview/"+preview.Filename, previewFile, preview.Size, preview.Header.Get("Content-Type"))
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

func (s *articleService) GetAllArticles() ([]model.Article, error) {
	articles, err := s.repo.GetAllPublished()

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *articleService) PublishArticle(articleId string, userID uint) error {
	err := s.repo.ChangeStatus(articleId, userID, true)
	return err
}

func (s *articleService) DeleteArticle(articleId string, userID uint) error {
	err := s.repo.Delete(articleId, userID)
	return err
}

func (s *articleService) UpdateArticle(
	articleID string,
	userID uint,
	title, content string,
	preview *multipart.FileHeader,
	attachments []*multipart.FileHeader,
) (*model.Article, error) {
	article, err := s.repo.FindByID(articleID)
	if err != nil {
		return nil, fmt.Errorf("article not found: %w", err)
	}

	if article.UserID != userID {
		return nil, fmt.Errorf("article does not belong to user")
	}

	if preview != nil {
		file, err := preview.Open()
		if err != nil {
			return nil, fmt.Errorf("preview open failed: %w", err)
		}
		defer file.Close()

		previewURL, err := s.minio.UploadFile("articles", "preview/"+preview.Filename, file, preview.Size, preview.Header.Get("Content-Type"))
		if err != nil {
			return nil, fmt.Errorf("preview upload failed: %w", err)
		}

		article.PreviewURL = previewURL
	}

	article.Title = title
	article.Content = content

	if err := s.repo.DeleteAttachmentsByArticleID(article.ID); err != nil {
		return nil, fmt.Errorf("failed to delete old attachments: %w", err)
	}
	article.Attachments = nil
	for _, file := range attachments {
		if len(article.Attachments) >= 10 {
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

	if err := s.repo.Update(article); err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	return article, nil
}

func (a articleService) GetAllArticlesAfterTime(time time.Time) ([]model.Article, error) {
	articles, err := a.repo.GetAllPublishedAfterTime(time)

	if err != nil {
		return nil, err
	}

	return articles, nil
}
