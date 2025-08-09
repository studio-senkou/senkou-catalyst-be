package services

import (
	"senkou-catalyst-be/integrations/midtrans"
)

type PaymentMethodsService interface {
	GetAllAvailablePaymentMethods() ([]midtrans.PaymentMethodConfig, error)
	GetPaymentMethodsByType(paymentType string) ([]midtrans.PaymentMethodConfig, error)
	GetPaymentMethodByChannelCode(channelCode string) (*midtrans.PaymentMethodConfig, error)
	GetPaymentMethodTypes() ([]string, error)
}

type PaymentMethodsServiceInstance struct{}

func NewPaymentMethodsService() PaymentMethodsService {
	return &PaymentMethodsServiceInstance{}
}

func (s *PaymentMethodsServiceInstance) GetAllAvailablePaymentMethods() ([]midtrans.PaymentMethodConfig, error) {
	return midtrans.GetPaymentMethods()
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodsByType(paymentType string) ([]midtrans.PaymentMethodConfig, error) {
	return midtrans.GetPaymentMethodsByType(paymentType)
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodByChannelCode(channel string) (*midtrans.PaymentMethodConfig, error) {
	return midtrans.GetPaymentMethodByChannel(channel)
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodTypes() ([]string, error) {
	methods, err := s.GetAllAvailablePaymentMethods()
	if err != nil {
		return nil, err
	}

	typeMap := make(map[string]bool)
	for _, method := range methods {
		typeMap[method.PaymentType] = true
	}

	var types []string
	for paymentType := range typeMap {
		types = append(types, paymentType)
	}

	return types, nil
}
