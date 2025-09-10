package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
)

type SubscriptionOrderService interface {
	CreateNewSubscriptionOrder(userID uint32, subID uint32, request *dtos.CreateSubscriptionOrderDTO) *errors.CustomError
	GetOrderByUserAndSubscription(orderID uint32, userID uint32) (*models.SubscriptionOrder, *errors.CustomError)
	UpdateSubscriptionOrder(orderID string, request *dtos.UpdateSubscriptionOrderDTO) *errors.CustomError
}

type SubscriptionOrderServiceInstance struct {
	SubscriptionOrderRepository  repositories.SubscriptionOrderRepository
	PaymentTransactionRepository repositories.PaymentTransactionRepository
}

func NewSubscriptionOrderService(
	subscriptionOrderRepository repositories.SubscriptionOrderRepository,
) SubscriptionOrderService {
	return &SubscriptionOrderServiceInstance{
		SubscriptionOrderRepository: subscriptionOrderRepository,
	}
}

func (s *SubscriptionOrderServiceInstance) CreateNewSubscriptionOrder(userID uint32, subID uint32, request *dtos.CreateSubscriptionOrderDTO) *errors.CustomError {

	orderID := uuid.New()

	newOrder := &models.SubscriptionOrder{
		ID:             orderID,
		UserID:         userID,
		SubscriptionID: subID,
	}

	if err := s.SubscriptionOrderRepository.StoreNewSubscriptionOrder(newOrder); err != nil {
		return errors.Internal("Failed to create subscription order", err.Error())
	}

	return nil
}

func (s *SubscriptionOrderServiceInstance) GetOrderByUserAndSubscription(orderID uint32, userID uint32) (*models.SubscriptionOrder, *errors.CustomError) {
	subscriptionOrder, err := s.SubscriptionOrderRepository.FindOrderByUserAndSubscription(userID, orderID)
	if err != nil {
		return nil, errors.NotFound("Subscription order not found")
	}

	return subscriptionOrder, nil
}

func (s *SubscriptionOrderServiceInstance) UpdateSubscriptionOrder(orderID string, request *dtos.UpdateSubscriptionOrderDTO) *errors.CustomError {

	return nil
}
