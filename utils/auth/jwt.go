package auth

import (
	"errors"
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret []byte
}

func NewJWTManager(secret string) (*JWTManager, error) {
	if len(secret) == 0 {
		return nil, errors.New("JWT secret cannot be empty")
	}
	return &JWTManager{Secret: []byte(secret)}, nil
}

func (j *JWTManager) GenerateToken(payload any, expiry time.Time) (*dtos.GeneratedToken, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     expiry.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &dtos.GeneratedToken{
		Token:     signedToken,
		ExpiresAt: fmt.Sprintf("%d", expiry.Unix()),
	}, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.Secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
