package dtos

type BaseMidtransNotification struct {
	Currency          string `json:"currency"`
	CustomField1      string `json:"custom_field1"`
	ExpiryTime        string `json:"expiry_time"`
	FraudStatus       string `json:"fraud_status"`
	GrossAmount       string `json:"gross_amount"`
	MerchantID        string `json:"merchant_id"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	SettlementTime    string `json:"settlement_time"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	StatusMessage     string `json:"status_message"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`
	TransactionTime   string `json:"transaction_time"`
}

type BankVA struct {
	Bank        string  `json:"bank"`
	VANumber    *string `json:"va_number"`
	GrossAmount *string `json:"gross_amount"`
}

type BankTransferNotification struct {
	BaseMidtransNotification

	PaymentAmount []BankVA `json:"payment_amounts"`
	VANumbers     []BankVA `json:"va_numbers"`
}

type UpdateSubscriptionOrderDTO struct {
	OrderID     string `json:"order_id" validate:"required"`
	FraudStatus string `json:"fraud_status" validate:"required"`

	// Transaction details
	TransactionID     *string `json:"transaction_id"`
	TransactionStatus *string `json:"transaction_status"`
	TransactionTime   *string `json:"transaction_time"`

	// Metadata
	SignatureKey   *string `json:"signature_key"`
	SettlementTime *string `json:"settlement_time"`
	ExpiryTime     *string `json:"expiry_time"`
}
