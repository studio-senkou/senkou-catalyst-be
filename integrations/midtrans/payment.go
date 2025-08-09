package midtrans

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

//go:embed constants/payment-methods.yaml
var paymentMethodsFS embed.FS

type PaymentMethodConfig struct {
	Name        string  `yaml:"name"          json:"name"`
	PaymentType string  `yaml:"payment_type"  json:"payment_type"`
	Channel     string  `yaml:"payment_channel" json:"channel"`
	MinAmount   float64 `yaml:"min_amount"    json:"min_amount"`
	MaxAmount   float64 `yaml:"max_amount"    json:"max_amount"`
	Description string  `yaml:"description"   json:"description"`
	LogoURL     string  `yaml:"logo_url"      json:"logo_url"`
}

type PaymentMethodsConfig struct {
	PaymentMethods []PaymentMethodConfig `yaml:"payment_methods"`
}

var paymentMethodsConfig *PaymentMethodsConfig

func GetPaymentMethods() ([]PaymentMethodConfig, error) {
	if paymentMethodsConfig == nil {
		if err := loadPaymentMethods(); err != nil {
			return nil, err
		}
	}

	return paymentMethodsConfig.PaymentMethods, nil
}

func GetPaymentMethodsByType(paymentType string) ([]PaymentMethodConfig, error) {
	allMethods, err := GetPaymentMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	var filteredMethods []PaymentMethodConfig
	for _, method := range allMethods {
		if method.PaymentType == paymentType {
			filteredMethods = append(filteredMethods, method)
		}
	}

	if len(filteredMethods) == 0 {
		return nil, fmt.Errorf("no payment methods found for type: %s", paymentType)
	}

	return filteredMethods, nil
}

func GetPaymentMethodByChannel(paymentChannel string) (*PaymentMethodConfig, error) {
	allMethods, err := GetPaymentMethods()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment methods: %w", err)
	}

	for _, method := range allMethods {
		if method.Channel == paymentChannel {
			return &method, nil
		}
	}

	return nil, fmt.Errorf("payment method not found for channel: %s", paymentChannel)
}

func loadPaymentMethods() error {
	data, err := paymentMethodsFS.ReadFile("constants/payment-methods.yaml")
	if err != nil {
		return fmt.Errorf("failed to read payment methods file: %w", err)
	}

	paymentMethodsConfig = &PaymentMethodsConfig{}
	if err := yaml.Unmarshal(data, paymentMethodsConfig); err != nil {
		return fmt.Errorf("failed to unmarshal payment methods: %w", err)
	}

	return nil
}

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSettled  PaymentStatus = "settled"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusCanceled PaymentStatus = "canceled"
	PaymentStatusDenied   PaymentStatus = "denied"
	PaymentStatusExpired  PaymentStatus = "expired"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

func ParseStatus(status string) PaymentStatus {
	switch status {
	case "pending":
		return PaymentStatusPending
	case "settlement":
		return PaymentStatusSettled
	case "deny":
		return PaymentStatusDenied
	case "expire":
		return PaymentStatusExpired
	case "cancel":
		return PaymentStatusCanceled
	case "refund":
		return PaymentStatusRefunded
	default:
		return PaymentStatusFailed
	}
}
