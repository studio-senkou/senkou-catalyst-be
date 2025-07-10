package controllers

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PredefinedCategoryController struct {
	pdService services.PredefinedCategoryService
}

func NewPredefinedCategoryController(pdService services.PredefinedCategoryService) *PredefinedCategoryController {
	return &PredefinedCategoryController{
		pdService: pdService,
	}
}

// Create a new predefined category
// @Summary Create a new predefined category
// @Description Create a new predefined category
// @Tags Predefined Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param CreatePDCategoryDTO body dtos.CreatePDCategoryDTO true "Create Predefined Category DTO"
// @Success 201 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /predefined-categories [post]
func (h *PredefinedCategoryController) StoreCategory(c *fiber.Ctx) error {
	categoryPDDTO := new(dtos.CreatePDCategoryDTO)

	if err := utils.Validate(c, categoryPDDTO); err != nil {
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

	if existingCategory, _ := h.pdService.GetPredefinedCategoryByName(categoryPDDTO.Name); existingCategory != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Predefined category already exists",
			"error":   "You cannot create a predefined category with the same name",
		})
	}

	if err := h.pdService.StoreCategory(categoryPDDTO); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create predefined category",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Predefined category created successfully",
	})
}

// Get all predefined categories
// @Summary Get all predefined categories
// @Description Retrieve all predefined categories
// @Tags Predefined Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PredefinedCategory
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /predefined-categories [get]
func (h *PredefinedCategoryController) GetPredefinedCategories(c *fiber.Ctx) error {
	predefinedCategories, err := h.pdService.GetAllPredefinedCategories()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve predefined categories",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"predefined_categories": predefinedCategories,
		},
		"message": "Predefined categories retrieved successfully",
	})
}

// Update a predefined category
// @Summary Update a predefined category
// @Description Update a predefined category by ID
// @Tags Predefined Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pcID path int true "Predefined Category ID"
// @Param UpdatePDCategoryDTO body dtos.UpdatePDCategoryDTO true "Update Predefined Category DTO"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 404 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /predefined-categories/{pcID} [put]
func (h *PredefinedCategoryController) UpdatePredefinedCategory(c *fiber.Ctx) error {
	pcID := c.Params("pcID")

	if pcID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Predefined Category ID is required",
			"error":   "Invalid request parameters",
		})
	}

	parsedPCID, err := strconv.ParseUint(pcID, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Predefined Category ID",
			"error":   err.Error(),
		})
	}

	updatePDCategoryDTO := new(dtos.UpdatePDCategoryDTO)

	if err := utils.Validate(c, updatePDCategoryDTO); err != nil {
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

	if err := h.pdService.UpdatePredefinedCategory(updatePDCategoryDTO, uint32(parsedPCID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update predefined category",
			"error":   err.Error(),
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully update predefined category",
	})
}

// Delete a predefined category
// @Summary Delete a predefined category
// @Description Delete a predefined category by ID
// @Tags Predefined Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pcID path int true "Predefined Category ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 404 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /predefined-categories/{pcID} [delete]
func (h *PredefinedCategoryController) DeletePredefinedCategory(c *fiber.Ctx) error {
	pcID := c.Params("pcID")

	if pcID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Predefined Category ID is required",
			"error":   "Invalid request parameters",
		})
	}

	parsedPCID, err := strconv.ParseUint(pcID, 10, 32)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Predefined Category ID",
			"error":   err.Error(),
		})
	}

	// Check if the predefined category exists before attempting to delete
	if _, err := h.pdService.GetPredefinedCategoryByName(pcID); err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Predefined category not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve predefined category",
			"error":   err.Error(),
		})

	}
	
	if err := h.pdService.DeletePredefinedCategory(uint32(parsedPCID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete predefined category",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted predefined category",
	})
}