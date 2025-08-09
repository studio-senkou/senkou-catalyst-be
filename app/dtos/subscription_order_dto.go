package dtos

type CreateSubscriptionOrderDTO struct {
	PaymentType    string  `json:"payment_type" validate:"required"`
	PaymentChannel string  `json:"payment_channel" validate:"required"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
}

func (dto *CreateSubscriptionOrderDTO) ErrorMessages() map[string]string {
	return map[string]string{
		"PaymentType.required":    "Payment type is required",
		"PaymentChannel.required": "Payment channel is required",
		"Amount.required":         "Amount is required",
		"Amount.gt":               "Amount must be greater than 0",
	}
}
