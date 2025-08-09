package repositories

import (
	"senkou-catalyst-be/app/models"

	"gorm.io/gorm"
)

type PaymentTransactionRepository interface {
	CreateTransaction(transaction *models.PaymentTransaction) error
	FindByID(transactionID string) (*models.PaymentTransaction, error)
	Update(transaction *models.PaymentTransaction) error
}

type PaymentTransactionRepositoryInstance struct {
	DB *gorm.DB
}

func NewPaymentTransactionRepository(db *gorm.DB) PaymentTransactionRepository {
	return &PaymentTransactionRepositoryInstance{
		DB: db,
	}
}

func (r *PaymentTransactionRepositoryInstance) CreateTransaction(transaction *models.PaymentTransaction) error {
	if err := r.DB.Create(transaction).Error; err != nil {
		return err
	}

	return nil
}

func (r *PaymentTransactionRepositoryInstance) FindByID(transactionID string) (*models.PaymentTransaction, error) {
	var transaction models.PaymentTransaction
	if err := r.DB.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *PaymentTransactionRepositoryInstance) Update(transaction *models.PaymentTransaction) error {
	if err := r.DB.Save(transaction).Error; err != nil {
		return err
	}

	return nil
}
