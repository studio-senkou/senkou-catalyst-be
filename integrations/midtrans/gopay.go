package midtrans

import "github.com/midtrans/midtrans-go/coreapi"

func setGopayParams(chargeReq *coreapi.ChargeReq) error {

	chargeReq.Gopay = &coreapi.GopayDetails{
		EnableCallback: true,
	}

	return nil
}
