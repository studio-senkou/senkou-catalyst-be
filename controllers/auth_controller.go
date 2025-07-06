package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		authService,
		userService,
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

	if err := utils.Validate(c, loginRequestDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  vErr.Errors,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	userID, err := h.userService.VerifyCredentials(loginRequestDTO.Email, loginRequestDTO.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	accessToken, refreshToken, err := h.authService.GenerateToken(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
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

	if err := utils.Validate(c, refreshTokenRequestDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Validation failed",
				"errors":  vErr.Errors,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	userID, err := h.authService.ValidateRefreshToken(refreshTokenRequestDTO.RefreshToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	accessToken, refreshToken, err := h.authService.GenerateToken(userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
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
	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user ID",
			"error":   err.Error(),
		})
	}

	if err := h.authService.InvalidateSession(uint(userID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to logout",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
