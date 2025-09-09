package services

import (
	stderr "errors"
	"fmt"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/auth"
	"senkou-catalyst-be/utils/config"
	"senkou-catalyst-be/utils/mailer"
	"senkou-catalyst-be/utils/query"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Activate(token string) *errors.AppError
	Create(user *models.User, merchant *models.Merchant) (*models.User, *errors.AppError)
	GetAll(params *query.QueryParams) (*[]models.User, *query.PaginationResponse, *errors.AppError)
	GetUserDetail(userID uint32) (*models.User, *errors.AppError)
	IsEmailVerified(userID uint32) (bool, *errors.AppError)
	SendEmailActivation(user *models.User) *errors.AppError
	VerifyCredentials(email, password string) (uint32, *errors.AppError)
	VerifyIsAnAdministrator(userID uint32) (bool, *errors.AppError)
}

type UserServiceInstance struct {
	UserRepository            repositories.UserRepository
	EmailActivationRepository repositories.EmailActivationRepository
	MerchantRepository        repositories.MerchantRepository
}

func NewUserService(userRepository repositories.UserRepository, merchantRepository repositories.MerchantRepository, emailActivationRepo repositories.EmailActivationRepository) UserService {
	return &UserServiceInstance{
		UserRepository:            userRepository,
		EmailActivationRepository: emailActivationRepo,
		MerchantRepository:        merchantRepository,
	}
}

// Create a new user in the database
// Returns the created user or an error if any
func (s *UserServiceInstance) Create(user *models.User, merchant *models.Merchant) (*models.User, *errors.AppError) {
	hashedPassword, err := user.HashPassword()

	if err != nil {
		return nil, errors.NewAppError(500, "Failed to hash password")
	}

	user.Password = hashedPassword

	if _, err := s.UserRepository.FindByEmail(user.Email); err != nil {
		if !stderr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewAppError(400, "User already exists")
		}
	}

	createdUser, err := s.UserRepository.Create(user)
	if err != nil {
		return nil, errors.NewAppError(500, "Failed to create user")
	}

	// If merchant is nil, skip creating merchant
	// This asume that the user is not a merchant
	// Otherwise, create the merchant and link it to the user
	if merchant == nil {
		return createdUser, nil
	}

	merchant.OwnerID = createdUser.ID

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

// Verify user credentials during login
// Returns the user ID if credentials are valid, or an error if invalid
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

// Check if the user's email is verified
// Returns true if the email is verified, false otherwise, or an error if any
func (s *UserServiceInstance) IsEmailVerified(userID uint32) (bool, *errors.AppError) {
	user, err := s.UserRepository.FindByID(userID)

	if err != nil {
		return false, errors.NewAppError(500, "Failed to find user by ID")
	}

	emailVerified := user.MustVerifyEmail()

	return emailVerified, nil
}

// Send email activation to the user
// Returns an error if any
func (s *UserServiceInstance) SendEmailActivation(user *models.User) *errors.AppError {
	secret := config.MustGetEnv("AUTH_SECRET")
	tokenManager, err := auth.NewJWTManager(secret)
	if err != nil {
		return errors.NewAppError(500, "Failed to create token manager")
	}

	verificationClaims := map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"type":    "account-activation",
	}
	verificationToken, err := tokenManager.GenerateToken(verificationClaims, time.Now().Add(24*time.Hour))
	if err != nil {
		return errors.NewAppError(500, "Failed to generate verification token")
	}

	mail, err := mailer.NewMailFromTemplate(user.Email, "Catalyst - Account activation", "account-activation.html", map[string]any{
		"ActivationLink": config.MustGetEnv("APP_FE_URL") + "/verify?token=" + verificationToken.Token,
	})
	if err != nil {
		return errors.NewAppError(500, "Failed to create email")
	}

	if err := mail.Send(); err != nil {
		return errors.NewAppError(500, "Failed to send email")
	}

	tokenExpiresAtUnix, err := strconv.ParseInt(verificationToken.ExpiresAt, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return errors.NewAppError(500, "Failed to parse expiration time")
	}

	tokenExpiresAt := time.Unix(tokenExpiresAtUnix, 0)

	activation := &models.EmailActivationToken{
		UserID:    user.ID,
		Token:     verificationToken.Token,
		ExpiresAt: tokenExpiresAt,
	}

	if _, err := s.EmailActivationRepository.Create(activation); err != nil {
		return errors.NewAppError(500, "Failed to store email activation token")
	}

	return nil
}

// Activate user account using the provided token
// Returns an error if any
func (s *UserServiceInstance) Activate(token string) *errors.AppError {

	activation, err := s.EmailActivationRepository.FindByToken(token)
	if err != nil {
		return errors.NewAppError(400, "Invalid or expired activation token")
	} else if activation == nil {
		return errors.NewAppError(404, "Activation token not found")
	}

	if time.Now().After(activation.ExpiresAt) || activation.UsedAt != nil {
		return errors.NewAppError(400, "Activation token has expired")
	}

	// Validate token claims to ensure it matches the activation record
	secret := config.MustGetEnv("AUTH_SECRET")
	tokenManager, err := auth.NewJWTManager(secret)
	if err != nil {
		return errors.NewAppError(500, "Failed to create token manager")
	}

	claims, err := tokenManager.ValidateToken(token)
	if err != nil {
		return errors.NewAppError(400, "Invalid activation token")
	}

	payload := claims["payload"].(map[string]interface{})

	if payload["type"] != "account-activation" || uint32(payload["user_id"].(float64)) != activation.UserID {
		return errors.NewAppError(400, "Invalid activation token")
	}

	now := time.Now()
	activation.UsedAt = &now

	if _, err := s.EmailActivationRepository.Update(activation); err != nil {
		return errors.NewAppError(500, "Failed to update email activation token")
	}

	if _, err := s.UserRepository.Update(&models.User{
		ID:              activation.UserID,
		EmailVerifiedAt: &now,
	}); err != nil {
		return errors.NewAppError(500, "Failed to activate user account")
	}

	return nil
}
