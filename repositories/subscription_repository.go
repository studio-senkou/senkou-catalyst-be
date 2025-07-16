package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{
		DB: db,
	}
}

// Store a new subscription
// This function saves a new subscription to the database
// It returns an error if the subscription could not be saved
func (r *SubscriptionRepository) StoreNewSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	if err := r.DB.Create(subscription).Error; err != nil {
		return nil, err
	}

	return subscription, nil
}
