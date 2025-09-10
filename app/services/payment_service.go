package services

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/integrations/midtrans"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"
	"time"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) (any, *errors.CustomError)
	UpdatePayment(transactionID string, updateData *models.PaymentTransaction) *errors.CustomError
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

func (s *PaymentServiceInstance) CreatePayment(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) (any, *errors.CustomError) {
	orderID := uuid.New().String()

	if err := s.validateInputs(user, pm, amount); err != nil {
		return nil, err
	}

	chargeReq, err := s.Builder.BuildChargeRequest(user, *pm, amount, orderID)
	if err != nil {
		return nil, errors.Internal("Failed to build charge request", err.Error())
	}

	transactionID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.BadRequest("Invalid order ID format", err.Error())
	}

	coreAPI := s.Client.GetCoreAPIClient()
	chargeResp, chargeRequestError := coreAPI.ChargeTransaction(chargeReq)
	if chargeRequestError != nil {
		return nil, errors.Internal("Failed to charge transaction", chargeRequestError.Error())
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
		return nil, errors.Internal("Failed to create payment transaction", err.Error())
	}

	return chargeResp, nil
}

func (s *PaymentServiceInstance) validateInputs(user *models.User, pm *midtrans.PaymentMethodConfig, amount float64) *errors.CustomError {
	if user == nil {
		return errors.BadRequest("User is required", nil)
	}
	if user.Name == "" || user.Email == "" || user.Phone == "" {
		return errors.BadRequest("User name, email, and phone are required", nil)
	}
	if pm == nil {
		return errors.BadRequest("Payment method is required", nil)
	}
	if amount <= 0 {
		return errors.BadRequest("Amount must be greater than 0", nil)
	}
	return nil
}

func (s *PaymentServiceInstance) UpdatePayment(transactionID string, updateData *models.PaymentTransaction) *errors.CustomError {
	if updateData == nil {
		return errors.BadRequest("Update data is required", nil)
	}

	existingTransaction, err := s.TransactionRepository.FindByID(transactionID)
	if err != nil {
		return errors.NotFound("Transaction not found")
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

	if err := s.TransactionRepository.Update(existingTransaction); err != nil {
		return errors.Internal("Failed to update transaction", err.Error())
	}

	return nil
}
