package services

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/integrations/midtrans"
	"senkou-catalyst-be/repositories"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) (any, error)
	UpdatePayment(transactionID string, updateData *models.PaymentTransaction) error
}

type PaymentServiceInstance struct {
	Client                *midtrans.MidtransClient
	Builder               *midtrans.PaymentBuilder
	TransactionRepository repositories.PaymentTransactionRepository
}

func NewPaymentService(client *midtrans.MidtransClient, transactionRepository repositories.PaymentTransactionRepository) PaymentService {
	return &PaymentServiceInstance{
		Client:                client,
		Builder:               midtrans.NewPaymentBuilder(client),
		TransactionRepository: transactionRepository,
	}
}

func (s *PaymentServiceInstance) CreatePayment(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) (any, error) {
	orderID := uuid.New().String()

	if err := s.validateInputs(user, pm, amount); err != nil {
		return nil, err
	}

	chargeReq, err := s.Builder.BuildChargeRequest(user, *pm, amount, orderID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to build charge request: "+err.Error())
	}

	transactionID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid order ID format: "+err.Error())
	}

	coreAPI := s.Client.GetCoreAPIClient()
	chargeResp, chargeRequestError := coreAPI.ChargeTransaction(chargeReq)
	if chargeRequestError != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to charge transaction: "+chargeRequestError.Error())
	}

	var transactionTimePtr *time.Time
	if chargeResp.TransactionTime != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", chargeResp.TransactionTime)
		if err == nil {
			transactionTimePtr = &parsedTime
		}
	}

	var expiredAtPtr *time.Time
	if chargeResp.ExpiryTime != "" {
		parsedExpiry, err := time.Parse("2006-01-02 15:04:05", chargeResp.ExpiryTime)
		if err == nil {
			expiredAtPtr = &parsedExpiry
		}
	}

	transaction := &models.PaymentTransaction{
		ID:              transactionID,
		Amount:          amount,
		Currency:        chargeResp.Currency,
		ExpiredAt:       expiredAtPtr,
		FraudStatus:     chargeResp.FraudStatus,
		PaymentChannel:  pm.Channel,
		PaymentType:     pm.PaymentType,
		Status:          chargeResp.TransactionStatus,
		TransactionID:   &chargeResp.TransactionID,
		TransactionTime: transactionTimePtr,
	}

	if err := s.TransactionRepository.CreateTransaction(transaction); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create payment transaction: "+err.Error())
	}

	return chargeResp, nil
}

func (s *PaymentServiceInstance) validateInputs(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) error {
	if user == nil {
		return fiber.ErrBadRequest
	}
	if user.Name == "" || user.Email == "" || user.Phone == "" {
		return fiber.ErrBadRequest
	}
	if pm == nil {
		return fiber.ErrBadRequest
	}
	if amount <= 0 {
		return fiber.ErrBadRequest
	}
	return nil
}

func (s *PaymentServiceInstance) UpdatePayment(transactionID string, updateData *models.PaymentTransaction) error {
	if updateData == nil {
		return fiber.ErrBadRequest
	}

	existingTransaction, err := s.TransactionRepository.FindByID(transactionID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction not found: "+err.Error())
	}

	if updateData.Status != "" {
		existingTransaction.Status = updateData.Status
	}
	if updateData.FraudStatus != "" {
		existingTransaction.FraudStatus = updateData.FraudStatus
	}
	if updateData.TransactionID != nil {
		existingTransaction.TransactionID = updateData.TransactionID
	}
	if updateData.TransactionTime != nil {
		existingTransaction.TransactionTime = updateData.TransactionTime
	}
	if updateData.SignatureKey != nil {
		existingTransaction.SignatureKey = updateData.SignatureKey
	}
	if updateData.ExpiredAt != nil {
		existingTransaction.ExpiredAt = updateData.ExpiredAt
	}
	if updateData.SettledAt != nil {
		existingTransaction.SettledAt = updateData.SettledAt
	}

	return s.TransactionRepository.Update(existingTransaction)
}
