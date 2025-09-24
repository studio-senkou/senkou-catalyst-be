package controllers

import (
	"encoding/hex"
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	goth "senkou-catalyst-be/integrations/goth"
	"senkou-catalyst-be/utils/config"
	"senkou-catalyst-be/utils/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/o1egl/paseto"
)

type OAuthController struct {
	UserService services.UserService
	AuthService services.AuthService
}

func NewOAuthController(userService services.UserService, authService services.AuthService) *OAuthController {
	return &OAuthController{
		UserService: userService,
		AuthService: authService,
	}
}

// GoogleCallback handles the callback from Google OAuth
// @Summary Google OAuth Callback
// @Description Handles the callback from Google OAuth and creates a new user if not exists
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map{data=models.User}
// @Failure 400 {object} fiber.Map{}
// @Failure 500 {object} fiber.Map{}
// @Router /auth/google/callback [get]
func (c *OAuthController) GoogleCallback(ctx *fiber.Ctx) error {

	user, err := goth.CompleteUserAuth(ctx)
	if err != nil {
		return response.InternalError(ctx, "Failed to authenticate user", err)
	}

	existingUser, existingError := c.UserService.GetByEmail(user.Email)
	if existingError != nil {
		if existingError.Code != fiber.StatusNotFound {
			return response.InternalError(ctx, "Failed to get user by email", existingError)
		}
	}
	if existingUser != nil {
		return loginAndRedirect(ctx, c.AuthService, existingUser)
	}

	username := strings.Split(user.Email, "@")[0]

	createUserModel := &models.User{
		Name:     user.FirstName + " " + user.LastName,
		Email:    user.Email,
		Phone:    "",
		Password: []byte(""), // No need for password as it's OAuth
		Role:     "user",
		IsOauth:  true,
	}

	createMerchantModel := &models.Merchant{
		Name:     username + "'s Merchant",
		Username: username,
	}

	oauthAccount := &dtos.CreateOAuthAccountDTO{
		Provider:     "google",
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    user.ExpiresAt,
	}

	account, createError := c.UserService.CreateOAuth(createUserModel, oauthAccount, createMerchantModel)
	if createError != nil {
		return response.InternalError(ctx, "Failed to create user", createError)
	}

	return loginAndRedirect(ctx, c.AuthService, account)
}

func loginAndRedirect(ctx *fiber.Ctx, authService services.AuthService, user *models.User) error {
	if redirectURI := config.MustGetEnv("GOOGLE_REDIRECT_URL"); redirectURI != "" {

		appKeyStr := config.MustGetEnv("APP_KEY")
		appKey, err := hex.DecodeString(appKeyStr)
		if err != nil {
			return response.InternalError(ctx, "Failed to decode app key", err)
		}

		if len(appKey) != 32 {
			return response.InternalError(ctx, "Failed to generate user session", fmt.Sprintf("Invalid key length: got %d bytes, need 32 bytes", len(appKey)))
		}

		accessToken, refreshToken, appError := authService.GenerateToken(user.ID)
		if appError != nil {
			return response.InternalError(ctx, "Failed to generate user session", nil)
		}

		now := time.Now()
		payload := fiber.Map{
			"token":         accessToken,
			"refresh_token": refreshToken,
			"iss":           "catalyst-backend",
			"sub":           user.ID,
			"aud":           "catalyst-frontend",
			"iat":           now.Unix(),
			"exp":           now.Add(1 * time.Hour).Unix(),
		}

		token, err := paseto.NewV2().Encrypt(appKey, payload, nil)
		if err != nil {
			return response.InternalError(ctx, "Failed to encrypt token", err)
		}

		return ctx.Status(fiber.StatusPermanentRedirect).Redirect(
			fmt.Sprintf("%s?token=%s", redirectURI, token),
		)
	}

	errorPageURI := config.MustGetEnv("APP_FE_URL")

	return ctx.Status(fiber.StatusPermanentRedirect).Redirect(
		fmt.Sprintf("%s?error=%s", errorPageURI, "Failed to create user"),
	)
}
