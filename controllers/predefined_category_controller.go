package controllers

import (
	"senkou-catalyst-be/services"

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

func (h *PredefinedCategoryController) StoreCategory(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Predefined category created successfully",
	})
}
