package goth

import (
	"senkou-catalyst-be/integrations/goth/google"
	"senkou-catalyst-be/utils/config"

	"github.com/markbates/goth"
)

func InitOAuthProviders() {

	goth.UseProviders(
		google.NewGoogleOAuthBuilder(
			config.MustGetEnv("GOOGLE_CLIENT_KEY"),
			config.MustGetEnv("GOOGLE_CLIENT_SECRET"),
			config.MustGetEnv("GOOGLE_CALLBACK_URL"),
		),
	)

}
