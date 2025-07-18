package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	StoreSession(userID uint32, token string) error
	FindSessionByToken(token string) (*models.UserHasToken, error)
	DeleteUserSession(userID uint32) error
}

type AuthRepositoryInstance struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryInstance{DB: db}
}

// Store a new user session
// This function will store a new user session in the database
// It returns an error if the session could not be stored
func (r *AuthRepositoryInstance) StoreSession(userID uint32, token string) error {
	session := models.UserHasToken{
		UserID: userID,
		Token:  token,
	}

	return r.DB.Create(&session).Error
}

// Find a user session by token
// This function will find a user session by token in the database
// It returns the session if found, or an error if not found
func (r *AuthRepositoryInstance) FindSessionByToken(token string) (*models.UserHasToken, error) {
	var session models.UserHasToken
	err := r.DB.Where("token = ?", token).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// Delete a user session by user
// This function will delete a user session by user in the database
// It returns an error if the session could not be deleted
func (r *AuthRepositoryInstance) DeleteUserSession(userID uint32) error {
	return r.DB.Where("user_id = ?", userID).Delete(&models.UserHasToken{}).Error
}
