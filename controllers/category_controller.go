package controllers

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"

	"github.com/gofiber/fiber/v2"
)

type CategoryController interface {
	CreateCategory(c *fiber.Ctx) error
	GetCategories(c *fiber.Ctx) error
}

type categoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) CategoryController {
	return &categoryController{
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
func (h *categoryController) CreateCategory(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")

	if merchantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant ID is required",
		})
	}

	createCategoryDTO := new(dtos.CreateCategoryDTO)

	if err := utils.Validate(c, createCategoryDTO); err != nil {
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

	if category, _ := h.categoryService.GetCategoryByName(createCategoryDTO.Name, merchantID); category != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Category already exists",
		})
	}

	category, err := h.categoryService.CreateNewCategory(createCategoryDTO, merchantID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create category",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"category": category,
		},
		"message": "Category created successfully",
	})
}

func (h *categoryController) GetCategories(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")

	if merchantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant ID is required",
		})
	}

	categories, err := h.categoryService.GetAllCategoriesByMerchantID(merchantID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve categories",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Categories retrieved successfully",
		"data": fiber.Map{
			"categories": categories,
		},
	})
}

func (h *categoryController) UpdateCategory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Update category functionality is not implemented yet",
	})
}

func (h *categoryController) DeleteCategory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Delete category functionality is not implemented yet",
	})
}
