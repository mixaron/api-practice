package model

import (
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title       string `gorm:"size:255"`
	Content     string `gorm:"type:text"`
	PreviewURL  string
	IsPublished bool `gorm:"default:false"`
	UserID      uint
	User        User         `gorm:"foreignKey:UserID"`
	Attachments []Attachment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ArticleID"`
}
