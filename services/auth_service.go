package services

import (
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils"
	"time"
)

type AuthService interface {
	GenerateToken(userID uint) (string, string, error)
}

type authService struct {
	authRepository *repositories.AuthRepository
}

func NewAuthService(authRepository *repositories.AuthRepository) AuthService {
	return &authService{authRepository}
}

func (s *authService) GenerateToken(userID uint) (string, string, error) {
	token, err := utils.GenerateToken(userID, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(userID, time.Now().Add(30*24*time.Hour))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
