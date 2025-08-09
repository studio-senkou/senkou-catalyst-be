package services

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	appError "senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
)

type SubscriptionOrderService interface {
	CreateNewSubscriptionOrder(userID uint32, subID uint32, request *dtos.CreateSubscriptionOrderDTO) *appError.AppError
	GetOrderByUserAndSubscription(orderID uint32, userID uint32) (*models.SubscriptionOrder, *appError.AppError)
	UpdateSubscriptionOrder(orderID string, request *dtos.UpdateSubscriptionOrderDTO) *appError.AppError
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

func (s *SubscriptionOrderServiceInstance) CreateNewSubscriptionOrder(userID uint32, subID uint32, request *dtos.CreateSubscriptionOrderDTO) *appError.AppError {

	orderID := uuid.New()

	newOrder := &models.SubscriptionOrder{
		ID:             orderID,
		UserID:         userID,
		SubscriptionID: subID,
	}

	if err := s.SubscriptionOrderRepository.StoreNewSubscriptionOrder(newOrder); err != nil {
		return appError.NewAppError(500, "Failed to create subscription order: "+err.Error())
	}

	return nil
}

func (s *SubscriptionOrderServiceInstance) GetOrderByUserAndSubscription(orderID uint32, userID uint32) (*models.SubscriptionOrder, *appError.AppError) {
	subscriptionOrder, err := s.SubscriptionOrderRepository.FindOrderByUserAndSubscription(userID, orderID)
	if err != nil {
		return nil, appError.NewAppError(404, "Subscription order not found: "+err.Error())
	}

	return subscriptionOrder, nil
}

func (s *SubscriptionOrderServiceInstance) UpdateSubscriptionOrder(orderID string, request *dtos.UpdateSubscriptionOrderDTO) *appError.AppError {

	return nil
}
