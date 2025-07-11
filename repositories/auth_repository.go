package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	StoreSession(userID uint, token string) error
	FindSessionByToken(token string) (*models.UserHasToken, error)
	DeleteUserSession(userID uint) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) StoreSession(userID uint, token string) error {
	session := models.UserHasToken{
		UserID: userID,
		Token:  token,
	}

	return r.db.Create(&session).Error
}

func (r *authRepository) FindSessionByToken(token string) (*models.UserHasToken, error) {
	var session models.UserHasToken
	err := r.db.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *authRepository) DeleteUserSession(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserHasToken{}).Error
}
