package repository

import (
	"api-practice/internal/model"
	"fmt"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *model.Article) error
	GetAllPublished() ([]model.Article, error)
	ChangeStatus(articleID string, userID uint, b bool) error
	Delete(articleID string, UserID uint) error
	Update(article *model.Article) error
	FindByID(id string) (*model.Article, error)
	DeleteAttachmentsByArticleID(articleID uint) error
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepo{db}
}

func (r *articleRepo) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepo) GetAllPublished() ([]model.Article, error) {
	var articles []model.Article
	if err := r.db.
		Preload("User").
		Preload("Attachments").
		Where("is_published = ?", true).
		Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *articleRepo) ChangeStatus(id string, userID uint, b bool) error {
	var article model.Article
	errFind := r.db.First(&article, id).Error
	if errFind != nil {
		return fmt.Errorf("article not found")
	}
	if article.UserID != userID {
		return fmt.Errorf("article does not belong to user")
	}

	if article.IsPublished == true {
		return fmt.Errorf("article already published")
	}
	r.db.Model(&article).Update("is_published", b)

	return nil
}

func (r *articleRepo) Delete(id string, userID uint) error {
	var article model.Article
	errFind := r.db.First(&article, id).Error
	if errFind != nil {
		return fmt.Errorf("article not found")
	}
	if article.UserID != userID {
		return fmt.Errorf("article does not belong to user")
	}

	r.db.Unscoped().Delete(&article)

	return nil
}

func (r *articleRepo) Update(article *model.Article) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(article).Error
}

func (r *articleRepo) FindByID(id string) (*model.Article, error) {
	var article model.Article
	if err := r.db.Preload("Attachments").First(&article, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepo) DeleteAttachmentsByArticleID(articleID uint) error {
	return r.db.Where("article_id = ?", articleID).Unscoped().Delete(&model.Attachment{}).Error
}
