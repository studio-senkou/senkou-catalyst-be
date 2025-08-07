package midtrans

type PaymentBuilder struct {
	client *MidtransClient
}

func NewPaymentBuilder(client *MidtransClient) *PaymentBuilder {
	return &PaymentBuilder{
		client: client,
	}
}
