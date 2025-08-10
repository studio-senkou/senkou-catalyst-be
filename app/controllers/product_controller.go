package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/storage"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	UserService    services.UserService
	ProductService services.ProductService
}

func NewProductController(productService services.ProductService, userService services.UserService) *ProductController {
	return &ProductController{
		ProductService: productService,
		UserService:    userService,
	}
}

// Create a new affilition product
// @Summary Create a new product
// @Description Create a new product for a merchant
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body dtos.CreateProductDTO true "Product data"
// @Success 201 {object} fiber.Map{message=string,data=fiber.Map{product=models.Product}}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products [post]
func (h *ProductController) CreateProduct(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("userID"))
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return response.InternalError(c, "Failed to parse user ID", fmt.Sprintf("Invalid user ID: %v", err.Error()))
	}

	createProductDTO := new(dtos.CreateProductDTO)

	if validationErrors, err := validator.ValidateFormData(c, createProductDTO); err != nil {
		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	} else if len(validationErrors) > 0 {
		return response.BadRequest(c, "Validation failed", validationErrors)
	}

	productPhoto, err := c.FormFile("photo")
	if err != nil {
		return response.BadRequest(c, "Failed to retrieve product photo", err.Error())
	} else if productPhoto == nil {
		return response.BadRequest(c, "Product photo is required", nil)
	}

	user, userErr := h.UserService.GetUserDetail(uint32(userID))
	if userErr != nil {
		return response.InternalError(c, "Failed to retrieve user details", userErr.Details)
	}

	if len(user.Merchants) == 0 {
		return response.BadRequest(c, "Cannot create product", "User does not have any associated merchants")
	}

	productPhotoPath, err := storage.UploadFileToStorage(productPhoto, "products", "PD", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "fail",
			"message": "Failed to upload product photo",
			"error":   err.Error(),
		})
	}

	createProductDTO.Photo = productPhotoPath

	createdProduct, appError := h.ProductService.CreateProduct(createProductDTO, user.Merchants[0].ID)

	if appError != nil {
		return response.InternalError(c, "Failed to create product", appError.Details)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"data": fiber.Map{
			"product": createdProduct,
		},
	})
}

// Get all products
// @Summary Get all products
// @Description Retrieve all products from the database
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fiber.Map{message=string,data=fiber.Map{products=[]models.Product}}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products [get]
func (h *ProductController) GetAllProducts(c *fiber.Ctx) error {
	products, appError := h.ProductService.GetAllProducts()

	if appError != nil {
		return response.InternalError(c, "Cannot retrieve products", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"data": fiber.Map{
			"products": products,
		},
	})
}

// Get products by merchant ID
// @Summary Get products by merchant ID
// @Description Retrieve all products associated with a specific merchant ID
// @Tags Products
// @Accept json
// @Produce json
// @Param merchantID path string true "Merchant ID"
// @Success 200 {object} fiber.Map{message=string,data=fiber.Map{products=[]models.Product}}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /merchants/{merchantID}/products [get]
func (h *ProductController) GetProductByMerchant(c *fiber.Ctx) error {
	merchantID := c.Params("merchantID")

	if merchantID == "" {
		return response.BadRequest(c, "Cannot continue to retrieve products", "Merchant ID is required")
	}

	products, appError := h.ProductService.GetProductsByMerchantID(merchantID)

	if appError != nil {
		switch appError.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "No products found for the specified merchant")
		case fiber.StatusInternalServerError:
			return response.InternalError(c, "Failed to retrieve products", appError.Details)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"data": fiber.Map{
			"products": products,
		},
	})
}

// Get a product by its ID
// @Summary Get a product by ID
// @Description Retrieve a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param productID path string true "Product ID"
// @Success 200 {object} fiber.Map{message=string,data=fiber.Map{product=models.Product}}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID} [get]
func (h *ProductController) GetProductByID(c *fiber.Ctx) error {
	productID := c.Params("productID")

	if productID == "" {
		return response.BadRequest(c, "Cannot continue to retrieve product information", "Product ID is required")
	}

	product, appError := h.ProductService.GetProductByID(productID)

	if appError != nil {
		switch appError.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Product not found")
		case fiber.StatusInternalServerError:
			return response.InternalError(c, "Failed to retrieve product", appError.Details)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product retrieved successfully",
		"data": fiber.Map{
			"product": product,
		},
	})
}

// Update a product
// @Summary Update a product
// @Description Update an existing product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Param product body dtos.UpdateProductDTO true "Updated product data"
// @Success 200 {object} fiber.Map{message=string,data=fiber.Map{product=models.Product}}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID} [put]
func (h *ProductController) UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("productID")

	if productID == "" {
		return response.BadRequest(c, "Cannot continue to update product", "Product ID is required")
	}

	updatedProductDTO := new(dtos.UpdateProductDTO)

	if err := validator.Validate(c, updatedProductDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	updatedProduct, appError := h.ProductService.UpdateProduct(updatedProductDTO, productID)

	if appError != nil {
		switch appError.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Product not found")
		case fiber.StatusInternalServerError:
			return response.InternalError(c, "Failed to update product", appError.Details)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"data": fiber.Map{
			"product": updatedProduct,
		},
	})
}

// Delete a product
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID} [delete]
func (h *ProductController) DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("productID")

	if productID == "" {
		return response.BadRequest(c, "Product ID is required", nil)
	}

	if err := h.ProductService.DeleteProduct(productID); err != nil {
		switch err.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Product not found")
		case fiber.StatusInternalServerError:
			return response.InternalError(c, "Failed to delete product", err.Details)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
