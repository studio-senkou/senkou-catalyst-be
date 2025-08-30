package services

import (
	stderr "errors"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/query"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user models.User) (*models.User, *errors.AppError)
	GetAll(params *query.QueryParams) (*[]models.User, *query.PaginationResponse, *errors.AppError)
	GetUserDetail(userID uint32) (*models.User, *errors.AppError)
	VerifyCredentials(email, password string) (uint32, *errors.AppError)
	VerifyIsAnAdministrator(userID uint32) (bool, *errors.AppError)
}

type UserServiceInstance struct {
	UserRepository     repositories.UserRepository
	MerchantRepository repositories.MerchantRepository
}

func NewUserService(userRepository repositories.UserRepository, merchantRepository repositories.MerchantRepository) UserService {
	return &UserServiceInstance{
		UserRepository:     userRepository,
		MerchantRepository: merchantRepository,
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

	if user, err := s.UserRepository.FindByEmail(user.Email); err != nil || user != nil {
		if !stderr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(400, "User already exists")
		}
	}

	createdUser, err := s.UserRepository.Create(&user)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to create user")
	}

	merchant := &models.Merchant{
		ID:      "",
		Name:    user.Name + " Merchant",
		OwnerID: createdUser.ID,
	}

	uuidStr := uuid.New().String()

	merchant.ID = strings.ReplaceAll(uuidStr, "-", "")[:16]

	createdMerchant, err := s.MerchantRepository.Create(merchant)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to create merchant")
	}

	createdUser.Merchants = append(createdUser.Merchants, createdMerchant)

	return createdUser, nil
}

// Get all users from the database
// Returns a slice of User models or an error if the operation fails
func (s *UserServiceInstance) GetAll(params *query.QueryParams) (*[]models.User, *query.PaginationResponse, *errors.AppError) {
	users, total, err := s.UserRepository.FindAll(params)
	if err != nil {
		return nil, nil, errors.NewAppError(500, "Failed to retrieve users")
	}

	pagination := query.CalculatePagination(params.Page, params.Limit, total)

	return users, pagination, nil
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
