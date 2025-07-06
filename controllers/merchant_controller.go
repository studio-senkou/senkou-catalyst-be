package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MerchantController struct {
	merchantService services.MerchantService
}

func NewMerchantController(merchantService services.MerchantService) *MerchantController {
	return &MerchantController{
		merchantService: merchantService,
	}
}

// Create merchant account
// @Summary Create Merchant
// @Description Create a merchant account for the user
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dtos.CreateMerchantRequestDTO body dtos.CreateMerchantRequestDTO true "Create merchant request"
// @Success 200 {object} fiber.Map{message=string,data=models.Merchant}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants [post]
func (h *MerchantController) CreateMerchant(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user ID",
			"error":   err.Error(),
		})
	}

	merchants, err := h.merchantService.GetUserMerchants(uint(userID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	if merchants != nil && len(*merchants) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "You already have a merchant",
			"error":   "You can only have one merchant per user",
		})
	}

	createMerchantRequestDTO := new(dtos.CreateMerchantRequestDTO)

	if err := utils.Validate(c, createMerchantRequestDTO); err != nil {
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

	userMerchants, _ := h.merchantService.GetUserMerchants(uint(userID))

	if userMerchants != nil && len(*userMerchants) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "You already have a merchant",
			"error":   "You can only have one merchant per user",
		})
	}

	merchant, err := h.merchantService.CreateMerchant(createMerchantRequestDTO, uint(userID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create merchant",
			"error":   err.Error(),
		})
	}

	if merchant == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create merchant",
			"error":   "Merchant creation returned nil",
		})
	}

	return c.JSON(fiber.Map{
		"data":    merchant,
		"message": "Merchant created successfully",
	})
}

// Get user merchants
// @Summary Get User Merchants
// @Description Retrieve all merchants associated with the user
// @Tags Merchant
// @Security BearerAuth
// @Success 200 {object} fiber.Map{data=fiber.Map{merchants=[]models.Merchant},message=string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants [get]
func (h *MerchantController) GetUserMerchants(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse user ID",
			"error":   err.Error(),
		})
	}

	merchants, err := h.merchantService.GetUserMerchants(uint(userID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	if merchants == nil || len(*merchants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "You have no merchants",
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"merchants": merchants,
		},
		"message": "Merchants retrieved successfully",
	})
}
