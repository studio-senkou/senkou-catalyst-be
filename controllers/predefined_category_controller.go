package controllers

import (
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"

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
	fmt.Println("returning from role middleware")
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
