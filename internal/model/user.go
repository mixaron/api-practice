package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email        string    `json:"email" gorm:"unique" validate:"required,email"`
	Password     string    `json:"password" validate:"required,min=8"`
	LastOnlineAt time.Time `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.LastOnlineAt.IsZero() {
		u.LastOnlineAt = time.Now()
	}
	return
}
