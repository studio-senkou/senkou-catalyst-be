package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type OAuthRepository interface {
	Store(oauthAccount *models.OauthAccount) (*models.OauthAccount, error)
}

type OAuthRepositoryInstance struct {
	db *gorm.DB
}

func NewOAuthRepository(db *gorm.DB) OAuthRepository {
	return &OAuthRepositoryInstance{db}
}

func (r *OAuthRepositoryInstance) Store(oauthAccount *models.OauthAccount) (*models.OauthAccount, error) {
	if err := r.db.Create(oauthAccount).Error; err != nil {
		return nil, err
	}
	return oauthAccount, nil
}
