package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type SubscriptionOrderRepository interface {
	StoreNewSubscriptionOrder(order *models.SubscriptionOrder) error
	FindOrderByUserAndSubscription(userID uint32, subID uint32) (*models.SubscriptionOrder, error)
	FindByOrderID(orderID string) (*models.SubscriptionOrder, error)
	UpdateOrderTransaction(orderID string, order *models.SubscriptionOrder) error
}

type SubscriptionOrderRepositoryInstance struct {
	DB *gorm.DB
}

func NewSubscriptionOrderRepository(db *gorm.DB) SubscriptionOrderRepository {
	return &SubscriptionOrderRepositoryInstance{
		DB: db,
	}
}

// Store a new subscription order
func (r *SubscriptionOrderRepositoryInstance) StoreNewSubscriptionOrder(order *models.SubscriptionOrder) error {
	if err := r.DB.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func (r *SubscriptionOrderRepositoryInstance) FindOrderByUserAndSubscription(userID uint32, subID uint32) (*models.SubscriptionOrder, error) {
	order := new(models.SubscriptionOrder)

	if err := r.DB.Where("user_id = ? AND subscription_id = ?", userID, subID).First(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *SubscriptionOrderRepositoryInstance) FindByOrderID(orderID string) (*models.SubscriptionOrder, error) {
	order := new(models.SubscriptionOrder)

	if err := r.DB.Where("order_id = ?", orderID).First(order).Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (r *SubscriptionOrderRepositoryInstance) UpdateOrderTransaction(orderID string, order *models.SubscriptionOrder) error {
	if err := r.DB.Where("id = ?", orderID).Updates(order).Error; err != nil {
		return err
	}

	return nil
}
