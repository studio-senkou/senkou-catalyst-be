package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/errors"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils"
	"time"
)

type AuthService interface {
	GenerateToken(userID uint32) (*dtos.GeneratedToken, *dtos.GeneratedToken, *errors.AppError)
	ValidateRefreshToken(token string) (uint32, *errors.AppError)
	InvalidateSession(userID uint32) *errors.AppError
}

type AuthServiceInstance struct {
	AuthRepository repositories.AuthRepository
}

func NewAuthService(authRepository repositories.AuthRepository) AuthService {
	return &AuthServiceInstance{AuthRepository: authRepository}
}

// Generate token and refresh token for the user
// This function generates a JWT token and a refresh token for the user
// It stores the refresh token in the database for later validation
func (s *AuthServiceInstance) GenerateToken(userID uint32) (*dtos.GeneratedToken, *dtos.GeneratedToken, *errors.AppError) {
	token, err := utils.GenerateToken(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, nil, errors.NewAppError(500, "Failed to generate token")
	}

	refreshToken, err := utils.GenerateToken(userID, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, nil, errors.NewAppError(500, "Failed to generate refresh token")
	}

	if err := s.AuthRepository.StoreSession(userID, refreshToken.Token); err != nil {
		return nil, nil, errors.NewAppError(500, "Failed to store session")
	}

	return token, refreshToken, nil
}

// Validate the refresh token
// This function checks if the provided refresh token exists in the database
// It returns the userID if the token is valid, otherwise it returns an error
func (s *AuthServiceInstance) ValidateRefreshToken(token string) (uint32, *errors.AppError) {
	session, err := s.AuthRepository.FindSessionByToken(token)

	if err != nil {
		return 0, errors.NewAppError(500, "Failed to validate refresh token")
	}

	if session == nil {
		return 0, errors.NewAppError(401, "Invalid refresh token")
	}

	return session.UserID, nil
}

// Invalidate the session for the user
// This function deletes the user's session from the database
// It is used to log out the user by removing their refresh token
// If the session is successfully deleted, it returns nil, otherwise it returns an error
func (s *AuthServiceInstance) InvalidateSession(userID uint32) *errors.AppError {
	if err := s.AuthRepository.DeleteUserSession(userID); err != nil {
		return errors.NewAppError(500, "Failed to invalidate session")
	}

	return nil
}
