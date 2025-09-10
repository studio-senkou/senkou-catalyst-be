package services

import (
	"senkou-catalyst-be/integrations/midtrans"
	"senkou-catalyst-be/platform/errors"
)

type PaymentMethodsService interface {
	GetAllAvailablePaymentMethods() ([]midtrans.PaymentMethodConfig, *errors.CustomError)
	GetPaymentMethodsByType(paymentType string) ([]midtrans.PaymentMethodConfig, *errors.CustomError)
	GetPaymentMethodByChannelCode(channelCode string) (*midtrans.PaymentMethodConfig, *errors.CustomError)
	GetPaymentMethodTypes() ([]string, *errors.CustomError)
}

type PaymentMethodsServiceInstance struct{}

func NewPaymentMethodsService() PaymentMethodsService {
	return &PaymentMethodsServiceInstance{}
}

func (s *PaymentMethodsServiceInstance) GetAllAvailablePaymentMethods() ([]midtrans.PaymentMethodConfig, *errors.CustomError) {
	methods, err := midtrans.GetPaymentMethods()
	if err != nil {
		return nil, errors.Internal("Failed to get payment methods", err.Error())
	}
	return methods, nil
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodsByType(paymentType string) ([]midtrans.PaymentMethodConfig, *errors.CustomError) {
	methods, err := midtrans.GetPaymentMethodsByType(paymentType)
	if err != nil {
		return nil, errors.Internal("Failed to get payment methods by type", err.Error())
	}
	return methods, nil
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodByChannelCode(channel string) (*midtrans.PaymentMethodConfig, *errors.CustomError) {
	method, err := midtrans.GetPaymentMethodByChannel(channel)
	if err != nil {
		return nil, errors.Internal("Failed to get payment method by channel", err.Error())
	}
	return method, nil
}

func (s *PaymentMethodsServiceInstance) GetPaymentMethodTypes() ([]string, *errors.CustomError) {
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
