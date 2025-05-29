package model

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	URL       string
	FileName  string
	ArticleID uint
}
