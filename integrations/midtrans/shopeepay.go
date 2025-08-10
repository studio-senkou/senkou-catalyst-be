package midtrans

import "github.com/midtrans/midtrans-go/coreapi"

func setShopeePayParams(chargeReq *coreapi.ChargeReq) error {
	chargeReq.PaymentType = coreapi.PaymentTypeShopeepay

	chargeReq.ShopeePay = &coreapi.ShopeePayDetails{
		// EnableCallback: true,
		CallbackUrl: "https://google.com", // TODO: Replace with actual callback URL
	}

	return nil
}
