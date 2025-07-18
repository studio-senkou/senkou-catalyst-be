package services

import (
	"errors"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type UserService interface {
	Create(user models.User) (*models.User, error)
	GetAll() (*[]models.User, error)
	GetUserDetail(userID uint32) (*models.User, error)
	VerifyCredentials(email, password string) (uint32, error)
	VerifyIsAnAdministrator(userID uint32) (bool, error)
}

type userService struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

// Create a new user in the database
// Returns the created user or an error if any
func (s *userService) Create(user models.User) (*models.User, error) {
	hashedPassword, err := user.HashPassword()

	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	return s.UserRepository.Create(&user)
}

// Get all users from the database
// Returns a slice of User models or an error if the operation fails
func (s *userService) GetAll() (*[]models.User, error) {
	return s.UserRepository.FindAll()
}

func (s *userService) VerifyCredentials(email, password string) (uint32, error) {
	user, err := s.UserRepository.FindByEmail(email)

	if err != nil {
		return 0, err
	}

	if !user.CheckPassword(password) {
		return 0, errors.New("invalid email or password")
	}

	return user.ID, nil
}

// Get user detail by its ID
// Returns the user model or an error if any
func (s *userService) GetUserDetail(userID uint32) (*models.User, error) {
	user, err := s.UserRepository.FindByID(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Verify if the user is an administrator
// Returns true if the user is an administrator, false otherwise, or an error if any
func (s *userService) VerifyIsAnAdministrator(userID uint32) (bool, error) {
	user, err := s.UserRepository.FindByID(userID)

	if err != nil {
		return false, err
	}

	return user.Role == "admin", nil
}
