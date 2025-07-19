package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MerchantController struct {
	MerchantService services.MerchantService
}

func NewMerchantController(merchantService services.MerchantService) *MerchantController {
	return &MerchantController{
		MerchantService: merchantService,
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
	userID, err := strconv.ParseUint(userIDStr, 10, 32)

	if userID == 0 || err != nil {
		return throw.BadRequest(c, "Cannot continue to create merchant", "Failed to parse user ID")
	}

	merchants, appError := h.MerchantService.GetUserMerchants(uint32(userID))
	if appError != nil {
		return throw.InternalError(c, "Cannot continue to create merchant", appError.Details)
	}

	if len(merchants) > 0 {
		return throw.BadRequest(c, "Cannot continue to create merchant", "You can only have one merchant per user")
	}

	createMerchantRequestDTO := new(dtos.CreateMerchantRequestDTO)

	if err := utils.Validate(c, createMerchantRequestDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"request": createMerchantRequestDTO,
			"error":   err.Error(),
		})
	}

	userMerchants, _ := h.MerchantService.GetUserMerchants(uint32(userID))

	if len(userMerchants) > 0 {
		return throw.BadRequest(c, "Cannot continue to create merchant", "Only one merchant is allowed per user")
	}

	merchant, appError := h.MerchantService.CreateMerchant(createMerchantRequestDTO, uint32(userID))

	if appError != nil {
		return throw.InternalError(c, "Failed to create merchant", appError.Details)
	}

	if merchant == nil {
		return throw.InternalError(c, "Failed to create merchant", "Merchant creation returned nil")
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
		return throw.BadRequest(c, "Cannot continue to retrieve user merchants", "Failed to parse user ID")
	}

	merchants, appError := h.MerchantService.GetUserMerchants(uint32(userID))

	if appError != nil {
		return throw.InternalError(c, "Cannot continue to retrieve user merchants due to internal error", appError.Details)
	}

	if len(merchants) == 0 {
		return throw.NotFound(c, "No merchants found for the user")
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
		return throw.BadRequest(c, "Cannot continue to retrieve merchant information", "Invalid merchant ID")
	}

	merchant, appError := h.MerchantService.GetMerchantByID(merchantID)

	if appError != nil {
		return throw.InternalError(c, "Cannot continue to retrieve merchant information due to internal error", appError.Details)
	}

	if merchant == nil {
		return throw.NotFound(c, "No merchant found with the provided identifiers")
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
		return throw.BadRequest(c, "Cannot continue to update merchant due to missing merchant ID", "Invalid merchant ID")
	}

	updateMerchantRequestDTO := new(dtos.UpdateMerchantRequestDTO)

	if err := utils.Validate(c, updateMerchantRequestDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	merchant, appError := h.MerchantService.UpdateMerchantByID(merchantID, updateMerchantRequestDTO)

	if appError != nil {
		return throw.InternalError(c, "Cannot update merchant due to internal error", appError.Details)
	}

	if merchant == nil {
		return throw.InternalError(c, "Cannot update merchant due to internal error", "Merchant update returned nil")
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
		return throw.BadRequest(c, "Cannot continue to delete the merchant due to invalid merchant ID", "Invalid merchant ID")
	}

	merchant, _ := h.MerchantService.GetMerchantByID(merchantID)

	if merchant == nil {
		return throw.NotFound(c, "Cannot remove merchant, Your merchant was not found")
	}

	appError := h.MerchantService.DeleteMerchantByID(merchantID)

	if appError != nil {
		return throw.InternalError(c, "Cannot delete merchant due to internal error", appError.Details)
	}

	return c.JSON(fiber.Map{
		"message": "Merchant deleted successfully",
	})
}
