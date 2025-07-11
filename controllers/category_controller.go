package controllers

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) CategoryController {
	return CategoryController{
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
		return throw.BadRequest(c, "Merchant ID is required", map[string]any{
			"Merchant ID": merchantID + " not found",
		})
	}

	createCategoryDTO := new(dtos.CreateCategoryDTO)

	if err := utils.Validate(c, createCategoryDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	if category, _ := h.categoryService.GetCategoryByName(createCategoryDTO.Name, merchantID); category != nil {
		return throw.BadRequest(c, "Category already exists", map[string]any{
			"category": createCategoryDTO.Name,
		})
	}

	category, err := h.categoryService.CreateNewCategory(createCategoryDTO, merchantID)

	if err != nil {
		return throw.InternalError(c, "Failed to create category", map[string]any{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"category": category,
		},
		"message": "Category created successfully",
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
		return throw.BadRequest(c, "Merchant ID is required", map[string]any{
			"Merchant ID": merchantID + " not found",
		})
	}

	categories, err := h.categoryService.GetAllCategoriesByMerchantID(merchantID)

	if err != nil {
		return throw.InternalError(c, "Failed to retrieve categories", map[string]any{
			"error": err.Error(),
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant ID and Category ID are required",
		})
	}

	parsedCategoryID, err := utils.StrToUint(categoryID)

	if err != nil {
		return throw.BadRequest(c, "Invalid Category ID", map[string]any{
			"error": err.Error(),
		})
	}

	updateCategoryDTO := new(dtos.UpdateCategoryDTO)

	if err := utils.Validate(c, updateCategoryDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	updatedCategory, err := h.categoryService.UpdateCategory(&models.Category{
		ID:         uint32(parsedCategoryID),
		Name:       updateCategoryDTO.Name,
		MerchantID: merchantID,
	})

	if err != nil {
		return throw.InternalError(c, "Failed to update category", map[string]any{
			"error": err.Error(),
		})
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
		return throw.BadRequest(c, "Merchant ID and Category ID are required", map[string]any{
			"merchantID": merchantID,
			"categoryID": categoryID,
		})
	}

	parsedCategoryID, err := utils.StrToUint(categoryID)

	if err != nil {
		return throw.BadRequest(c, "Invalid Category ID", map[string]any{
			"error": err.Error(),
		})
	}

	if err := h.categoryService.DeleteCategory(uint32(parsedCategoryID)); err != nil {
		return throw.DatabaseError(c, "Failed to delete category", "DELETE")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}
