package services

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	GenerateToken(userID uint) (*dtos.GeneratedToken, *dtos.GeneratedToken, error)
	ValidateRefreshToken(token string) (uint, error)
	InvalidateSession(userID uint) error
}

type authService struct {
	authRepository repositories.AuthRepository
}

func NewAuthService(authRepository repositories.AuthRepository) AuthService {
	return &authService{authRepository}
}

func (s *authService) GenerateToken(userID uint) (*dtos.GeneratedToken, *dtos.GeneratedToken, error) {
	token, err := utils.GenerateToken(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := utils.GenerateToken(userID, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return nil, nil, err
	}

	if err := s.authRepository.StoreSession(userID, refreshToken.Token); err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
}

func (s *authService) ValidateRefreshToken(token string) (uint, error) {
	session, err := s.authRepository.FindSessionByToken(token)

	if err != nil {
		return 0, err
	}

	if session == nil {
		return 0, fiber.ErrNotFound
	}

	return session.UserID, nil
}

func (s *authService) InvalidateSession(userID uint) error {
	if err := s.authRepository.DeleteUserSession(userID); err != nil {
		return err
	}

	return nil
}
