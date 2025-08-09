package midtrans

import (
	"errors"
	"fmt"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/utils/config"

	mt "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type PaymentBuilder struct {
	client *MidtransClient
}

func NewPaymentBuilder(client *MidtransClient) *PaymentBuilder {
	return &PaymentBuilder{
		client: client,
	}
}

func (b *PaymentBuilder) BuildChargeRequest(user *models.User, pm PaymentMethodConfig, amount float64, orderID string) (*coreapi.ChargeReq, error) {
	grossAmount := int64(amount)

	chargeReq := &coreapi.ChargeReq{
		PaymentType: b.getPaymentType(pm.PaymentType),
		TransactionDetails: mt.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		CustomerDetails: &mt.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
			Phone: user.Phone,
		},
	}

	if notificationUrl := b.getNotificationURL(); notificationUrl != "" {
		chargeReq.CustomField1 = &notificationUrl
	}

	if err := b.setPaymentMethodParams(chargeReq, &pm); err != nil {
		return nil, fmt.Errorf("failed to set payment method params: %w", err)
	}

	return chargeReq, nil
}

func (b *PaymentBuilder) getPaymentType(paymentType string) coreapi.CoreapiPaymentType {
	switch paymentType {
	case "bank_transfer":
		return coreapi.PaymentTypeBankTransfer
	case "e_wallet_gopay":
		return coreapi.PaymentTypeGopay
	case "e_wallet_shopeepay":
		return coreapi.PaymentTypeShopeepay
	case "e_wallet_dana":
		return "dana"
	case "qris":
		return coreapi.PaymentTypeQris
	case "otc_alfamart":
		return "alfamart"
	}

	return ""
}

func (b *PaymentBuilder) setPaymentMethodParams(chargeReq *coreapi.ChargeReq, pm *PaymentMethodConfig) error {
	switch pm.PaymentType {
	case "bank_transfer":
		return setBankTransferParams(chargeReq, BankChannel(pm.Channel))
	}

	return errors.New("unsupported payment method: " + pm.PaymentType)
}

func (b *PaymentBuilder) getNotificationURL() string {
	baseURL := config.GetEnv("APP_URL", "http://localhost:8080")
	return fmt.Sprintf("%s/api/v1/payments/notifications", baseURL)
}
