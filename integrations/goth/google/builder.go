package google

import (
	oauth "github.com/markbates/goth/providers/google"
)

type GoogleOAuthBuilder struct {
	ClientKey    string
	ClientSecret string
	CallbackURL  string
}

// NewGoogleOAuthBuilder creates a new instance of GoogleOAuthBuilder
func NewGoogleOAuthBuilder(clientKey, clientSecret, callbackURL string) *oauth.Provider {
	return oauth.New(clientKey, clientSecret, callbackURL)
}

// SetClientKey sets the client key for the Google OAuth builder
func (b *GoogleOAuthBuilder) SetClientKey(key string) *GoogleOAuthBuilder {
	b.ClientKey = key
	return b
}

// SetClientSecret sets the client secret for the Google OAuth builder
func (b *GoogleOAuthBuilder) SetClientSecret(secret string) *GoogleOAuthBuilder {
	b.ClientSecret = secret
	return b
}

// SetCallbackURL sets the callback URL for the Google OAuth builder
func (b *GoogleOAuthBuilder) SetCallbackURL(url string) *GoogleOAuthBuilder {
	b.CallbackURL = url
	return b
}

// Build constructs the Google OAuth provider using the configured parameters
func (b *GoogleOAuthBuilder) Build() *oauth.Provider {

	if b.ClientKey == "" {
		panic("Google client key is required")
	}

	if b.ClientSecret == "" {
		panic("Google client secret is required")
	}

	if b.CallbackURL == "" {
		panic("Google callback URL is required")
	}

	return oauth.New(b.ClientKey, b.ClientSecret, b.CallbackURL)
}
