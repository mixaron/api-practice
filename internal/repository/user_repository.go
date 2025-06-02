package repository

import (
	"api-practice/internal/db"
	"api-practice/internal/model"
	"fmt"
	"time"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	Save(user *model.User) error
	GetById(id uint) (*model.User, error)
	UpdateUserLastTime(user *model.User) error
}

type userRepoImpl struct{}

func NewUserRepository() UserRepository {
	return &userRepoImpl{}
}

func (r *userRepoImpl) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("user not found by email")
	}

	return &user, nil
}

func (r *userRepoImpl) Save(user *model.User) error {
	err := db.DB.Create(user).Error
	if err != nil {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func (r *userRepoImpl) GetById(id uint) (*model.User, error) {
	var user model.User
	err := db.DB.Where("id = ?", &id).First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("user not found by id")
	}

	return &user, nil
}

func (r *userRepoImpl) UpdateUserLastTime(user *model.User) error {
	err := db.DB.Model(&user).Update("last_online_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}
