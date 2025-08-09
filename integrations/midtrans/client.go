package midtrans

import (
	"senkou-catalyst-be/utils/config"

	m "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransClient struct {
	CoreAPI     coreapi.Client
	ServerKey   string
	Environment m.EnvironmentType
}

func NewMidtransClient() *MidtransClient {
	serverKey := config.MustGetEnv("MIDTRANS_SERVER_KEY")
	env := config.GetEnv("MIDTRANS_ENVIRONMENT", "sandbox")

	var environment m.EnvironmentType

	if env == "production" {
		environment = m.Production
	} else {
		environment = m.Sandbox
	}

	coreAPIClient := new(coreapi.Client)
	coreAPIClient.New(serverKey, environment)

	return &MidtransClient{
		CoreAPI:     *coreAPIClient,
		ServerKey:   serverKey,
		Environment: environment,
	}
}

func (mc *MidtransClient) GetCoreAPIClient() *coreapi.Client {
	return &mc.CoreAPI
}

func (mc *MidtransClient) GetCurrentEnvironment() m.EnvironmentType {
	return mc.Environment
}
