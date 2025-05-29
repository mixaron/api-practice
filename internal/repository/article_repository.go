package repository

import (
	"api-practice/internal/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *model.Article) error
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
