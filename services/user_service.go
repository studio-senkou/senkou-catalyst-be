package services

import (
	"errors"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type UserService interface {
	Create(user models.User) (*models.User, error)
	GetAll() (*[]models.User, error)
	VerifyCredentials(email, password string) (uint, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Create(user models.User) (*models.User, error) {
	hashedPassword, err := user.HashPassword()

	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	return s.repo.Create(&user)
}

func (s *userService) GetAll() (*[]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) VerifyCredentials(email, password string) (uint, error) {
	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return 0, err
	}

	if !user.CheckPassword(password) {
		return 0, errors.New("invalid email or password")
	}

	return uint(user.ID), nil
}
