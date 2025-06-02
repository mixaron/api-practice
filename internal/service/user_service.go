package service

import (
	"api-practice/internal/model"
	"api-practice/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type UserService interface {
	Register(user *model.User) error
	Authenticate(email, password string) (*model.User, error)
	GetProfile(userID uint) (*model.User, error)
	IsUserExists(userID uint) bool
	GetUserLastOnlineTime(userID string) (time.Time, error)
	SetUserLastOnlineTime(id string) error
}

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Register(user *model.User) error {
	if user.Password == "" {
		return fmt.Errorf("empty password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error cache password: %w", err)
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
		return nil, fmt.Errorf("invalid password")
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

func (s *UserServiceImpl) IsUserExists(userID uint) bool {
	if _, err := s.repo.GetById(userID); err != nil {
		return false
	}
	return true
}

func (s *UserServiceImpl) GetUserLastOnlineTime(userID string) (time.Time, error) {
	user, err := s.getUserById(userID)
	if err != nil {
		return time.Now(), err
	}
	return user.LastOnlineAt, nil
}

func (s *UserServiceImpl) SetUserLastOnlineTime(userID string) error {
	user, err := s.getUserById(userID)
	if err != nil {
		return err
	}
	user.LastOnlineAt = time.Now()

	errUpdate := s.repo.UpdateUserLastTime(user)

	if errUpdate != nil {
		return errUpdate
	}

	return nil
}

func (s *UserServiceImpl) getUserById(userID string) (*model.User, error) {
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.repo.GetById(uint(id))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}
