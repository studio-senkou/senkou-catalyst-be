package controllers

import (
	"senkou-catalyst-be/utils/storage"

	"github.com/gofiber/fiber/v2"
)

type StorageController struct{}

func NewStorageController() *StorageController {
	return &StorageController{}
}

func (s *StorageController) GetFromStorage(c *fiber.Ctx) error {

	filename := c.Params("*")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid filename"})
	}

	file, err := storage.DownloadFileFromStorage(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to download file from storage",
			"error":   err.Error(),
		})
	}

	return c.Send(file)
}
