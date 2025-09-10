package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService services.AuthService
	UserService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

// Login User
// @Summary Login user
// @Version 1.0
// @Description Login user with email and password to receive access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.LoginRequestDTO true "Request to authenticate"
// @Success 200 {object} dtos.LoginResponseDTO "Login successful response"
// @Router /auth/login [post]
func (h *AuthController) Login(c *fiber.Ctx) error {
	loginRequestDTO := new(dtos.LoginRequestDTO)

	if err := validator.Validate(c, loginRequestDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	userID, err := h.UserService.VerifyCredentials(loginRequestDTO.Email, loginRequestDTO.Password)

	if err != nil {
		return response.BadRequest(c, "Invalid email or password", nil)
	}

	emailVerified, err := h.UserService.IsEmailVerified(userID)

	if err != nil {
		return response.InternalError(c, "Failed to verify email status", err.Details)
	} else if !emailVerified {
		return response.Forbidden(c, "Email not verified. Please verify your email to proceed.")
	}

	accessToken, refreshToken, appError := h.AuthService.GenerateToken(userID)

	if appError != nil {
		return response.InternalError(c, "Failed to generate token", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data": dtos.LoginResponseDTO{
			AccessToken:        accessToken.Token,
			AccessTokenExpiry:  accessToken.ExpiresAt,
			RefreshToken:       refreshToken.Token,
			RefreshTokenExpiry: refreshToken.ExpiresAt,
		},
	})
}

// Refresh user session
// @Summary Refresh access token
// @Version 1.0
// @Description Refresh access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.RefreshTokenRequestDTO true "Request to refresh access token"
// @Success 200 {object} dtos.LoginResponseDTO "Token refreshed successfully"
// @Failure 400 {object} map[string]string "Validation failed"
// @Router /auth/refresh [put]
func (h *AuthController) RefreshToken(c *fiber.Ctx) error {
	refreshTokenRequestDTO := new(dtos.RefreshTokenRequestDTO)

	if err := validator.Validate(c, refreshTokenRequestDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	userID, err := h.AuthService.ValidateRefreshToken(refreshTokenRequestDTO.RefreshToken)

	if userID == 0 || err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	userID, refreshError := h.AuthService.ValidateRefreshToken(refreshTokenRequestDTO.RefreshToken)

	if refreshError != nil {
		if refreshError.Code == 401 {
			return response.Unauthorized(c, fmt.Sprintf("Cannot continue to update your token because of %v", refreshError.Details))
		}

		return response.InternalError(c, "Failed to validate refresh token", map[string]any{
			"error": refreshError.Details,
		})
	}

	accessToken, refreshToken, tokenError := h.AuthService.GenerateToken(userID)

	if tokenError != nil {
		return response.InternalError(c, "Failed to generate token", tokenError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Token refreshed successfully",
		"data": dtos.LoginResponseDTO{
			AccessToken:        accessToken.Token,
			AccessTokenExpiry:  accessToken.ExpiresAt,
			RefreshToken:       refreshToken.Token,
			RefreshTokenExpiry: refreshToken.ExpiresAt,
		},
	})
}

// Logout user
// @Summary Logout user
// @Version 1.0
// @Description Logout user by invalidating their session
// @Tags Auth
// @Security BearerAuth
// @Success 200 {object} map[string]string "Logout successful response"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /auth/logout [delete]
func (h *AuthController) Logout(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if userID == 0 || err != nil {
		return response.BadRequest(c, "Cannot continue to logout user", "User ID is not valid")
	}

	if err := h.AuthService.InvalidateSession(uint32(userID)); err != nil {
		return response.InternalError(c, "Failed to logout user", map[string]any{
			"error": err.Details,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
