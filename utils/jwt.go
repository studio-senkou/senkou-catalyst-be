package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(GetEnv("AUTH_SECRET", ""))

	if len(jwtSecret) == 0 {
		panic("JWT secret is not set. Please set the AUTH_SECRET environment variable.")
	}
}

func GenerateToken(payload any, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     expiry.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
