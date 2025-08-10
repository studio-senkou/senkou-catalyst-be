package midtrans

import "github.com/midtrans/midtrans-go/coreapi"

func setQrisParams(chargeReq *coreapi.ChargeReq) error {
	chargeReq.PaymentType = "qris"

	chargeReq.CustomExpiry = &coreapi.CustomExpiry{
		Unit:           "minute",
		ExpiryDuration: 15,
	}

	return nil
}
