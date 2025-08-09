package midtrans

import "github.com/midtrans/midtrans-go/coreapi"

type BankChannel string

const (
	BankChannelBCA     BankChannel = "bca"
	BankChannelBNI     BankChannel = "bni"
	BankChannelBRI     BankChannel = "bri"
	BankChannelMandiri BankChannel = "mandiri"
	BankChannelPermata BankChannel = "permata"
)

func setBankTransferParams(chargeReq *coreapi.ChargeReq, channel BankChannel) error {
	switch channel {
	case BankChannelBCA:
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: "bca",
		}
	}
	
	return nil
}