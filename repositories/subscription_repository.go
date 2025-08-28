package repositories

import (
	"senkou-catalyst-be/app/models"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	StoreNewSubscription(subscription *models.Subscription) (*models.Subscription, error)
	SubscribeUser(sub *models.UserSubscription) error
	FindAllSubscriptions() ([]*models.Subscription, error)
	FindByID(id uint32) (*models.Subscription, error)
	FindActiveSubscriptionByUserID(userID uint32) (*models.Subscription, error)
	FindFreeTierSubscription() (*models.Subscription, error)
	UpdateSubscription(updatedSubscription *models.Subscription) (*models.Subscription, error)
	DeleteSubscription(subscription *models.Subscription) error
	VerifyUserHasActiveSubscription(userID, subID uint32) (bool, error)
}

type SubscriptionRepositoryInstance struct {
	DB *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &SubscriptionRepositoryInstance{
		DB: db,
	}
}

// Store a new subscription
// This function saves a new subscription to the database
// It returns an error if the subscription could not be saved
func (r *SubscriptionRepositoryInstance) StoreNewSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	result := r.DB.Where("name = ?", subscription.Name).FirstOrCreate(subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

// Subscribe a user to a subscription
// This function store relation between user and subscription
// It returns an error if the subscription could not be created
func (r *SubscriptionRepositoryInstance) SubscribeUser(sub *models.UserSubscription) error {
	if err := r.DB.Create(sub).Error; err != nil {
		return err
	}

	return nil
}

// Find all the subscriptions
// This function retrieves all subscriptions from the database
// It returns a slice of subscriptions and an error if any
func (r *SubscriptionRepositoryInstance) FindAllSubscriptions() ([]*models.Subscription, error) {
	subscriptions := make([]*models.Subscription, 0)

	if err := r.DB.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// Find an active subscription by user id
// This function retrieves an active subscription for a specific user
// It returns the subscription and an error if any
func (r *SubscriptionRepositoryInstance) FindActiveSubscriptionByUserID(userID uint32) (*models.Subscription, error) {
	userSubscription := new(models.UserSubscription)
	subscription := new(models.Subscription)

	if err := r.DB.Where("user_id = ? AND is_active = ?", userID, true).First(userSubscription).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Plans").Find(subscription).Error; err != nil {
		return nil, err
	}

	return subscription, nil
}

// Find a subscription by its ID
// This function retrieves a subscription by its ID from the database
// It returns the subscription and an error if any
func (r *SubscriptionRepositoryInstance) FindByID(id uint32) (*models.Subscription, error) {
	subscription := new(models.Subscription)
	if err := r.DB.First(subscription, id).Error; err != nil {
		return nil, err
	}
	return subscription, nil
}

// Find a free tier subscription
// This function retrieve a subscription specific for free tier
// It returns the subscription and an error if any
func (r *SubscriptionRepositoryInstance) FindFreeTierSubscription() (*models.Subscription, error) {
	freeTierSubscription := new(models.Subscription)
	if err := r.DB.Where("name = $1 AND price = $2", "Free tier", 0).First(freeTierSubscription).Error; err != nil {
		return nil, err
	}
	return freeTierSubscription, nil
}

// Update an existing subscription
// This function updates an existing subscription in the database
// It returns the updated subscription and an error if any
func (r *SubscriptionRepositoryInstance) UpdateSubscription(updatedSubscription *models.Subscription) (*models.Subscription, error) {
	if err := r.DB.Save(updatedSubscription).Error; err != nil {
		return nil, err
	}

	return updatedSubscription, nil
}

// Delete a subscription
// This function deletes a subscription from the database
// It returns an error if the deletion fails
func (r *SubscriptionRepositoryInstance) DeleteSubscription(subscription *models.Subscription) error {
	if err := r.DB.Delete(subscription).Error; err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionRepositoryInstance) VerifyUserHasActiveSubscription(userID, subID uint32) (bool, error) {
	userSubscription := new(models.UserSubscription)

	err := r.DB.Where("user_id = ? AND sub_id = ?", userID, subID).First(userSubscription).Error

	if err != nil {
		return false, err
	}

	// Check if the user subscription is not expired
	if userSubscription.PaymentStatus == "declined" {
		return false, nil
	}

	// Check if the user has an active subscription but already expired
	if userSubscription.PaymentStatus == "settled" && userSubscription.IsActive && userSubscription.ExpiredAt.Before(time.Now()) {
		return false, nil
	}

	return true, nil
}
