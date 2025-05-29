package service

import (
	"api-practice/internal/model"
	"api-practice/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) error
	Authenticate(email, password string) (*model.User, error)
	GetProfile(id uint) (*model.User, error)
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(user *model.User) error {
	if user.Password == "" {
		return fmt.Errorf("Empty password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error cache password: %w", err)
	}

	user.Password = string(hashedPassword)

	return s.repo.Save(user)
}

func (s *UserServiceImpl) Authenticate(email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("Invalid password")
	}

	return user, nil
}

func (s *UserServiceImpl) GetProfile(userId uint) (*model.User, error) {
	user, err := s.repo.GetById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
