package controllers

import (
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/services"
	"senkou-catalyst-be/utils"
	"senkou-catalyst-be/utils/throw"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// Create a new affilition product
// @Summary Create a new product
// @Description Create a new product for a merchant
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param merchantID path string true "Merchant ID"
// @Param product body dtos.CreateProductDTO true "Product data"
// @Success 201 {object} fiber.Map{message=string,data=fiber.Map{product=models.Product}}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /merchants/{merchantID}/products [post]
func (h *ProductController) CreateProduct(c *fiber.Ctx) error {
	createProductDTO := new(dtos.CreateProductDTO)

	if err := utils.Validate(c, createProductDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	createdProduct, err := h.productService.CreateProduct(createProductDTO, createProductDTO.MerchantID)

	if err != nil {
		return throw.InternalError(c, "Failed to create product", map[string]any{
			"error": err.Error(),
		})
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
	products, err := h.productService.GetAllProducts()

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return throw.NotFound(c, "No products found")
		}

		return throw.InternalError(c, "Failed to retrieve products", map[string]any{
			"error": err.Error(),
		})
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
		return throw.BadRequest(c, "Merchant ID is required", nil)
	}

	products, err := h.productService.GetProductsByMerchantID(merchantID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return throw.NotFound(c, "No products found for this merchant")
		}

		return throw.InternalError(c, "Failed to retrieve products", map[string]any{
			"error": err.Error(),
		})
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
		return throw.BadRequest(c, "Product ID is required", nil)
	}

	product, err := h.productService.GetProductByID(productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return throw.NotFound(c, "Product not found")
		}

		return throw.InternalError(c, "Failed to retrieve product", map[string]any{
			"error": err.Error(),
		})
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
		return throw.BadRequest(c, "Product ID is required", nil)
	}

	updatedProductDTO := new(dtos.UpdateProductDTO)

	if err := utils.Validate(c, updatedProductDTO); err != nil {
		if vErr, ok := err.(*utils.ValidationError); ok {
			return throw.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return throw.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	updatedProduct, err := h.productService.UpdateProduct(updatedProductDTO, productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return throw.NotFound(c, "Product not found")
		}

		return throw.InternalError(c, "Failed to update product", map[string]any{
			"error": err.Error(),
		})
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
		return throw.BadRequest(c, "Product ID is required", nil)
	}

	if err := h.productService.DeleteProduct(productID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return throw.NotFound(c, "Product not found")
		}

		return throw.InternalError(c, "Failed to delete product", map[string]any{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}
