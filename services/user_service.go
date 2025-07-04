package services

import (
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type UserService interface {
	GetAll() ([]models.User, error)
	Create(user models.User) (models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) Create(user models.User) (models.User, error) {
	return s.repo.Create(user)
}
