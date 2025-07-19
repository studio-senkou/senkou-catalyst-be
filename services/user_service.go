package services

import (
	"senkou-catalyst-be/errors"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"
)

type UserService interface {
	Create(user models.User) (*models.User, *errors.AppError)
	GetAll() (*[]models.User, *errors.AppError)
	GetUserDetail(userID uint32) (*models.User, *errors.AppError)
	VerifyCredentials(email, password string) (uint32, *errors.AppError)
	VerifyIsAnAdministrator(userID uint32) (bool, *errors.AppError)
}

type UserServiceInstance struct {
	UserRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceInstance{
		UserRepository: userRepository,
	}
}

// Create a new user in the database
// Returns the created user or an error if any
func (s *UserServiceInstance) Create(user models.User) (*models.User, *errors.AppError) {
	hashedPassword, err := user.HashPassword()

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to hash password")
	}

	user.Password = hashedPassword

	createdUser, err := s.UserRepository.Create(&user)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to create user")
	}

	return createdUser, nil
}

// Get all users from the database
// Returns a slice of User models or an error if the operation fails
func (s *UserServiceInstance) GetAll() (*[]models.User, *errors.AppError) {
	users, err := s.UserRepository.FindAll()
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to retrieve users")
	}

	return users, nil
}

func (s *UserServiceInstance) VerifyCredentials(email, password string) (uint32, *errors.AppError) {
	user, err := s.UserRepository.FindByEmail(email)

	if err != nil {
		return 0, errors.NewAppError(500, "Failed to find user by email")
	}

	if !user.CheckPassword(password) {
		return 0, errors.NewAppError(400, "invalid email or password")
	}

	return user.ID, nil
}

// Get user detail by its ID
// Returns the user model or an error if any
func (s *UserServiceInstance) GetUserDetail(userID uint32) (*models.User, *errors.AppError) {
	user, err := s.UserRepository.FindByID(userID)

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to find user by ID")
	}

	return user, nil
}

// Verify if the user is an administrator
// Returns true if the user is an administrator, false otherwise, or an error if any
func (s *UserServiceInstance) VerifyIsAnAdministrator(userID uint32) (bool, *errors.AppError) {
	user, err := s.UserRepository.FindByID(userID)

	if err != nil {
		return false, errors.NewAppError(500, "Failed to find user by ID")
	}

	return user.Role == "admin", nil
}
