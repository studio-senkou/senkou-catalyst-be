package controllers

import (
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PredefinedCategoryController struct {
	PredefinedCategoryService services.PredefinedCategoryService
}

func NewPredefinedCategoryController(PDService services.PredefinedCategoryService) *PredefinedCategoryController {
	return &PredefinedCategoryController{
		PredefinedCategoryService: PDService,
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

	if err := validator.Validate(c, categoryPDDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if existingCategory, _ := h.PredefinedCategoryService.GetPredefinedCategoryByName(categoryPDDTO.Name); existingCategory != nil {
		return response.BadRequest(c, "Predefined category already exists", "Cannot create a predefined category with the same name")
	}

	predefinedCategory, appError := h.PredefinedCategoryService.StoreCategory(categoryPDDTO)
	if appError != nil {
		return response.InternalError(c, "Failed to create predefined category", appError.Details)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Predefined category created successfully",
		"data": fiber.Map{
			"predefined_category": predefinedCategory,
		},
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
	predefinedCategories, appError := h.PredefinedCategoryService.GetAllPredefinedCategories()

	if appError != nil {
		return response.InternalError(c, "Cannot retrieve predefined categories due to internal error", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Predefined categories retrieved successfully",
		"data": fiber.Map{
			"predefined_categories": predefinedCategories,
		},
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
	PCID := c.Params("pcID")

	if PCID == "" {
		return response.BadRequest(c, "Cannot continue updating predefined category", "Invalid predefined category ID")
	}

	parsedPCID, err := strconv.ParseUint(PCID, 10, 32)

	if err != nil {
		return response.BadRequest(c, "Invalid Predefined Category ID", map[string]any{
			"error": err.Error(),
		})
	}

	updatePDCategoryDTO := new(dtos.UpdatePDCategoryDTO)

	if err := validator.Validate(c, updatePDCategoryDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if appError := h.PredefinedCategoryService.UpdatePredefinedCategory(updatePDCategoryDTO, uint32(parsedPCID)); appError != nil {
		return response.InternalError(c, "Failed to update predefined category", appError.Details)
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
	PCID := c.Params("pcID")

	if PCID == "" {
		return response.BadRequest(c, "Cannot remove predefined category", "Invalid predefined category ID")
	}

	parsedPCID, err := strconv.ParseUint(PCID, 10, 32)

	if err != nil {
		return response.BadRequest(c, "Invalid Predefined Category ID", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.PredefinedCategoryService.DeletePredefinedCategory(uint32(parsedPCID)); err != nil {
		switch err.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Predefined category not found")
		case fiber.StatusInternalServerError:
			return response.InternalError(c, "Failed to delete predefined category", err.Details)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted predefined category",
	})
}
