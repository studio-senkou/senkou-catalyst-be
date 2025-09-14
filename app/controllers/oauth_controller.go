package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	goth "senkou-catalyst-be/integrations/goth"
	"senkou-catalyst-be/utils/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type OAuthController struct {
	UserService services.UserService
}

func NewOAuthController(userService services.UserService) *OAuthController {
	return &OAuthController{
		UserService: userService,
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
		fmt.Printf("User with email %s already exists", user.Email)
		return response.BadRequest(ctx, "User with this email already exists", nil)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    account,
	})
}
