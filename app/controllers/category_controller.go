package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/converter"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/validator"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// Create new category
// @Summary Create a new category
// @Description Create a new category for a merchant
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param merchantID path string true "Merchant ID"
// @Param CreateCategoryDTO body dtos.CreateCategoryDTO true "Create Category DTO"
// @Success 200 {object} fiber.Map{data=fiber.Map{category=models.Category}}
// @Failure 400 {object} fiber.Map{message=string, errors=[]string}
// @Failure 500 {object} fiber.Map{message=string, error=string}
// @Router /merchants/{merchantID}/categories [post]
func (h *CategoryController) CreateCategory(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")

	if merchantID == "" {
		return response.BadRequest(c, "Cannot continue to create category", "Invalid merchant ID")
	}

	createCategoryDTO := new(dtos.CreateCategoryDTO)

	if err := validator.Validate(c, createCategoryDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if category, err := h.categoryService.GetCategoryByName(createCategoryDTO.Name, merchantID); err != nil {
		if err.Code == 500 {
			return response.InternalError(c, "Cannot continue to create category due to internal error", err.Details)
		}
	} else if category != nil {
		return response.BadRequest(c, "Cannot continue to create category due to conflict", "Category already exists with the same name for this merchant")
	}

	category, appError := h.categoryService.CreateNewCategory(createCategoryDTO, merchantID)

	if appError != nil {
		return response.InternalError(c, "Cannot create the category errdue to internal error", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category created successfully",
		"data": fiber.Map{
			"category": category,
		},
	})
}

// Get all categories for a merchant
// @Summary Get all categories for a merchant
// @Description Retrieve all categories associated with a specific merchant
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fiber.Map{data=fiber.Map{categories=[]models.Category}}
// @Failure 400 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string, error=string}
// @Router /merchants/{merchantID}/categories [get]
func (h *CategoryController) GetCategories(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")

	if merchantID == "" {
		return response.BadRequest(c, "Cannot continue to retrieve categories", "Invalid merchant ID")
	}

	categories, appError := h.categoryService.GetAllCategoriesByMerchantID(merchantID)

	if appError != nil {
		return response.InternalError(c, "Cannot retrieve categories due to internal error", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Categories retrieved successfully",
		"data": fiber.Map{
			"categories": categories,
		},
	})
}

// Update category
// @Summary Update a category
// @Description Update an existing category for a merchant
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param merchantID path string true "Merchant ID"
// @Param categoryID path string true "Category ID"
// @Param UpdateCategoryDTO body dtos.UpdateCategoryDTO true "Update Category DTO"
// @Success 200 {object} fiber.Map{data=fiber.Map{category=models.Category}}
// @Failure 400 {object} fiber.Map{message=string,errors=[]string}
// @Failure 404 {object} fiber.Map{message=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants/{merchantID}/categories/{categoryID} [put]
func (h *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")
	categoryID := c.Params("categoryID")

	if merchantID == "" || categoryID == "" {
		return response.BadRequest(c, "Cannot continue to update category due to missing IDs", fmt.Sprintf("Invalid Merchant ID: %s, Invalid Category ID: %s", merchantID, categoryID))
	}

	parsedCategoryID, err := converter.StrToUint(categoryID)

	if err != nil {
		return response.BadRequest(c, "Cannot continue to update category due to invalid ID", "Invalid category ID: "+err.Error())
	}

	updateCategoryDTO := new(dtos.UpdateCategoryDTO)

	if err := validator.Validate(c, updateCategoryDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.ValidationError(c, "Validation failed", vErr.Errors)
		}

		return response.InternalError(c, "Internal server error", err.Error())
	}

	updatedCategory, appError := h.categoryService.UpdateCategory(&models.Category{
		ID:         uint32(parsedCategoryID),
		Name:       updateCategoryDTO.Name,
		MerchantID: merchantID,
	})

	if appError != nil {
		return response.InternalError(c, "Cannot update the category due to internal error", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"category": updatedCategory,
		},
		"message": "Category updated successfully",
	})
}

// Delete category
// @Summary Delete a category
// @Description Delete an existing category for a merchant
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param merchantID path string true "Merchant ID"
// @Param categoryID path string true "Category ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{message=string,error=string}
// @Failure 500 {object} fiber.Map{message=string,error=string}
// @Router /merchants/{merchantID}/categories/{categoryID} [delete]
func (h *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")
	categoryID := c.Params("categoryID")

	if merchantID == "" || categoryID == "" {
		return response.BadRequest(c, "Cannot continue to delete category due to missing IDs", fmt.Sprintf("Invalid Merchant ID: %s, Invalid Category ID: %s", merchantID, categoryID))
	}

	parsedCategoryID, err := converter.StrToUint(categoryID)

	if err != nil {
		return response.BadRequest(c, "Cannot continue to delete category due to invalid ID", "Invalid category ID: "+err.Error())
	}

	if err := h.categoryService.DeleteCategory(uint32(parsedCategoryID)); err != nil {
		return response.InternalError(c, "Cannot delete category due to internal error", err.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
