package controllers

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/services"
	"senkou-catalyst-be/utils/query"
	"senkou-catalyst-be/utils/response"
	"senkou-catalyst-be/utils/storage"
	"senkou-catalyst-be/utils/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	UserService    services.UserService
	ProductService services.ProductService
	ProductMetric  services.ProductInteractionService
}

func NewProductController(productService services.ProductService, userService services.UserService, productMetric services.ProductInteractionService) *ProductController {
	return &ProductController{
		ProductService: productService,
		UserService:    userService,
		ProductMetric:  productMetric,
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

	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, "Failed to parse multipart form data", err.Error())
	}

	photos := form.File["photos"]
	if len(photos) == 0 {
		return response.BadRequest(c, "At least one product photo required", nil)
	}

	user, userErr := h.UserService.GetUserDetail(uint32(userID))
	if userErr != nil {
		return response.InternalError(c, "Failed to retrieve user details", userErr.Details)
	}

	if len(user.Merchants) == 0 {
		return response.BadRequest(c, "Cannot create product", "User does not have any associated merchants")
	}

	var photoPaths []string
	for _, photo := range photos {
		if !storage.IsValidImageExtension(photo.Filename) {
			return response.BadRequest(c, "Invalid image format", fmt.Sprintf("File %s has an unsupported format", photo.Filename))
		}

		photoPath, err := storage.UploadFileToStorage(photo, "products", "PD", nil)
		if err != nil {
			return response.InternalError(c, "Failed to upload product photo", err.Error())
		}

		photoPaths = append(photoPaths, photoPath)
	}

	createProductDTO.Photos = photoPaths

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

// Upload a product photo
// @Summary Upload a product photo
// @Description Upload a photo for a specific product
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Param photo formData file true "Product photo"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID}/photos [post]
func (h *ProductController) UploadProductPhoto(c *fiber.Ctx) error {
	productId := c.Params("productID")
	if productId == "" {
		return response.BadRequest(c, "Product ID is required", nil)
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		return response.BadRequest(c, "Failed to parse photo", err.Error())
	}

	if !storage.IsValidImageExtension(photo.Filename) {
		return response.BadRequest(c, "Invalid image format", fmt.Sprintf("File %s has an unsupported format", photo.Filename))
	}

	// product, err := h.ProductService.GetProductByID(productId)
	product, appError := h.ProductService.GetProductByID(productId)
	if appError != nil {
		switch appError.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Product not found")
		default:
			return response.InternalError(c, "Failed to retrieve product", appError.Details)
		}
	} else if product == nil {
		return response.NotFound(c, "Product not found")
	} else if len(product.Photos) >= 5 {
		return response.BadRequest(c, "Cannot upload more photos", "Maximum of 5 photos allowed per product")
	}

	photoPath, err := storage.UploadFileToStorage(photo, "products", "PD", nil)
	if err != nil {
		return response.InternalError(c, "Failed to upload product photo", err.Error())
	}

	product.Photos.AddPhoto(photoPath)

	if err := h.ProductService.UpdateProductPhotos(product); err != nil {
		return response.InternalError(c, "Failed to update product", err.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully update product photos",
	})
}

// Delete product photo
// @Summary Delete the product photo by it's file path
// @Description Delete a photo for a specific product
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Param filePath path string true "File path"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID}/photos/{filePath} [delete]
func (h *ProductController) DeleteProductPhoto(c *fiber.Ctx) error {
	productId := c.Params("productID")
	if productId == "" {
		return response.BadRequest(c, "Product ID is required", nil)
	}

	filePath := c.Params("*")
	if filePath == "" {
		return response.BadRequest(c, "File path is required", nil)
	}

	product, err := h.ProductService.GetProductByID(productId)
	if err != nil {
		switch err.Code {
		case fiber.StatusNotFound:
			return response.NotFound(c, "Product not found")
		default:
			return response.InternalError(c, "Failed to retrieve product", err.Details)
		}
	}

	if err := storage.RemoveFileFromStorage(filePath); err != nil {
		return response.InternalError(c, "Failed to remove photo file", nil)
	}

	product.Photos.RemovePhoto(filePath)

	if err := h.ProductService.UpdateProductPhotos(product); err != nil {
		return response.InternalError(c, "Failed to update product", err.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully remove photo file",
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
	params := query.ParseQueryParams(c)

	products, pagination, appError := h.ProductService.GetAllProducts(params)

	if appError != nil {
		return response.InternalError(c, "Cannot retrieve products", appError.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Products retrieved successfully",
		"data": fiber.Map{
			"products":   products,
			"pagination": pagination,
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

// Send product log
// @Summary Sending product log for interaction metrics
// @Description Send a log for product interactions
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productID path string true "Product ID"
// @Param log body dtos.SendProductInteractionDTO true "Product interaction log"
// @Success 200 {object} fiber.Map{message=string}
// @Failure 400 {object} fiber.Map{error=string,details=any}
// @Failure 404 {object} fiber.Map{error=string,details=any}
// @Failure 500 {object} fiber.Map{error=string,details=any}
// @Router /products/{productID}/logs [post]
func (h *ProductController) SendProductLog(c *fiber.Ctx) error {
	productID := c.Params("productID")

	if productID == "" {
		return response.BadRequest(c, "Product ID is required", nil)
	}

	sendProductLogDTO := new(dtos.SendProductInteractionDTO)

	if err := validator.Validate(c, sendProductLogDTO); err != nil {
		if vErr, ok := err.(*validator.ValidationError); ok {
			return response.BadRequest(c, "Validation failed", map[string]any{
				"errors": vErr.Errors,
			})
		}

		return response.InternalError(c, "Internal server error", map[string]any{
			"error": err.Error(),
		})
	}

	parsedProductID := uuid.MustParse(productID)

	if err := h.ProductMetric.StoreLog(parsedProductID, sendProductLogDTO); err != nil {
		return response.InternalError(c, "Failed to store product interaction log", err.Details)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product interaction log sent successfully",
	})
}
