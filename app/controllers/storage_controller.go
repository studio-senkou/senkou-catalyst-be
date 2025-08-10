package controllers

import (
	"senkou-catalyst-be/utils/storage"

	"github.com/gofiber/fiber/v2"
)

type StorageController struct{}

func NewStorageController() *StorageController {
	return &StorageController{}
}

// GetFromStorage retrieves a file from the storage service
// @Summary Get a file from storage
// @Description Retrieve a file from the storage service by its filename
// @Tags Storage
// @Accept json
// @Produce json
// @Param filename path string true "Filename"
// @Success 200 {file} file "File content"
// @Failure 400 {object} fiber.Map{error=string}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /files/{filename} [get]
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
