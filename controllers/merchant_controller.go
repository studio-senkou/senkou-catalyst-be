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

// Get merchant by ID
// @Summary Get Merchant by ID
// @Description Retrieve a merchant by its ID
// @Tags Merchant
// @Security BearerAuth
// @Param id path string true "Merchant ID"
// @Success 200 {object} fiber.Map{data=models.Merchant,message=string}
// @Failure 400 {object} fiber.Map{message=string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants/{id} [get]
func (h *MerchantController) GetMerchantByID(c *fiber.Ctx) error {
	merchantID := c.Params("id")

	if merchantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid merchant ID",
		})
	}

	merchant, err := h.merchantService.GetMerchantByID(merchantID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	if merchant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Merchant not found",
		})
	}

	return c.JSON(fiber.Map{
		"data":    merchant,
		"message": "Merchant retrieved successfully",
	})
}

// Update merchant
// @Summary Update Merchant
// @Description Update a merchant's details
// @Tags Merchant
// @Security BearerAuth
// @Param id path string true "Merchant ID"
// @Param dtos.UpdateMerchantRequestDTO body dtos.UpdateMerchantRequestDTO true "Update Merchant request"
// @Success 200 {object} fiber.Map{data=fiber.Map{merchant=models.Merchant},message=string}
// @Failure 400 {object} fiber.Map{message=string,errors=[]string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants/{id} [put]
func (h *MerchantController) UpdateMerchant(c *fiber.Ctx) error {
	merchantID := c.Params("id")

	if merchantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid merchant ID",
		})
	}

	updateMerchantRequestDTO := new(dtos.UpdateMerchantRequestDTO)

	if err := utils.Validate(c, updateMerchantRequestDTO); err != nil {
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

	merchant, err := h.merchantService.UpdateMerchantByID(merchantID, updateMerchantRequestDTO)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update merchant",
			"error":   err.Error(),
		})
	}

	if merchant == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update merchant",
			"error":   "Merchant update returned nil",
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"merchant": merchant,
		},
		"message": "Merchant updated successfully",
	})
}

// Delete merchant
// @Summary Delete Merchant
// @Description Delete a merchant by its ID
// @Tags Merchant
// @Security BearerAuth
// @Param id path string true "Merchant ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants/{id} [delete]
func (h *MerchantController) DeleteMerchant(c *fiber.Ctx) error {
	merchantID := c.Params("id")

	if merchantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid merchant ID",
		})
	}

	merchant, _ := h.merchantService.GetMerchantByID(merchantID)

	if merchant == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Cannot remove merchant, Your merchant was not found",
		})
	}

	err := h.merchantService.DeleteMerchantByID(merchantID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete merchant",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Merchant deleted successfully",
	})
}
