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

// Find all the subscriptions
// This function retrieves all subscriptions from the database
// It returns a slice of subscriptions and an error if any
func (r *SubscriptionRepository) FindAllSubscriptions() ([]*models.Subscription, error) {
	subscriptions := make([]*models.Subscription, 0)

	if err := r.DB.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Find a subscription by its ID
// This function retrieves a subscription by its ID from the database
// It returns the subscription and an error if any
func (r *SubscriptionRepository) FindByID(id uint32) (*models.Subscription, error) {
	subscription := new(models.Subscription)
	if err := r.DB.First(subscription, id).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

// Update an existing subscription
// This function updates an existing subscription in the database
// It returns the updated subscription and an error if any
func (r *SubscriptionRepository) UpdateSubscription(updatedSubscription *models.Subscription) (*models.Subscription, error) {
	if err := r.DB.Save(updatedSubscription).Error; err != nil {
		return nil, err
	}

	return updatedSubscription, nil
}

// Delete a subscription
// This function deletes a subscription from the database
// It returns an error if the deletion fails
func (r *SubscriptionRepository) DeleteSubscription(subscription *models.Subscription) error {
	if err := r.DB.Delete(subscription).Error; err != nil {
		return err
	}

	return nil
}
