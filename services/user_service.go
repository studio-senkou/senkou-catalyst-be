package services

import (
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type UserService interface {
	GetAll() (*[]models.User, error)
	Create(user models.User) (*models.User, error)
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

