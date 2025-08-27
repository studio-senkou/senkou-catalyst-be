package services

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, *errors.AppError)
	GetProductByID(productID string) (*models.Product, *errors.AppError)
	GetProductsByMerchantID(merchantID string) ([]*models.Product, *errors.AppError)
	GetAllProducts() ([]*models.Product, *errors.AppError)
	UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, *errors.AppError)
	UpdateProductPhotos(product *models.Product) *errors.AppError
	DeleteProduct(productID string) *errors.AppError
	VerifyProductOwnership(productID string, userID uint32) *errors.AppError
}

type ProductServiceInstance struct {
	UserRepository    repositories.UserRepository
	ProductRepository repositories.ProductRepository
}

func NewProductService(productRepository repositories.ProductRepository, userRepository repositories.UserRepository) ProductService {
	return &ProductServiceInstance{
		ProductRepository: productRepository,
		UserRepository:    userRepository,
	}
}

// Create a new product
// This function will be used to create a new product via repository
// It returns the created product and an error if any
func (s *ProductServiceInstance) CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, *errors.AppError) {
	newProduct := &models.Product{
		ID:           uuid.New().String(),
		Title:        product.Title,
		Description:  product.Description,
		Price:        product.Price,
		Photos:       product.Photos,
		AffiliateURL: product.AffiliateURL,
		CategoryID:   product.CategoryID,
		MerchantID:   merchantID,
	}

	storedProduct, err := s.ProductRepository.StoreProduct(newProduct)

	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to create product: %v", err.Error()))
	}

	return storedProduct, nil
}

// Get a product by its ID
// This function retrieves a product from the repository by its ID
// It returns the product and an error if any
func (s *ProductServiceInstance) GetProductByID(productID string) (*models.Product, *errors.AppError) {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(404, "Product not found")
		}

		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to retrieve product: %v", err.Error()))
	}

	return product, nil
}

// Get products by merchant ID
// This function retrieves all products associated with a specific merchant ID
// It returns a slice of products and an error if any
func (s *ProductServiceInstance) GetProductsByMerchantID(merchantID string) ([]*models.Product, *errors.AppError) {
	products, err := s.ProductRepository.FindProductsByMerchantID(merchantID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError(404, "No products found for this merchant")
		}

		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to retrieve products: %v", err.Error()))
	}

	return products, nil
}

// Get all products
// This function retrieves all products from the repository
// It returns a slice of products and an error if any
func (s *ProductServiceInstance) GetAllProducts() ([]*models.Product, *errors.AppError) {
	products, err := s.ProductRepository.FindAllProducts()

	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to retrieve products: %v", err.Error()))
	}

	return products, nil
}

// Update a product
// This function updates an existing product in the repository
// It takes an updated product request and a product ID as parameters
// It returns the updated product and an error if any
func (s *ProductServiceInstance) UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, *errors.AppError) {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		return nil, errors.NewAppError(404, fmt.Sprintf("Product not found: %v", err.Error()))
	}

	if product == nil {
		// return nil, errors.New("product not found")
	}

	if updatedProduct.Title != nil {
		product.Title = *updatedProduct.Title
	}

	if updatedProduct.Description != nil {
		product.Description = *updatedProduct.Description
	}

	if updatedProduct.Price != nil {
		product.Price = *updatedProduct.Price
	}

	if updatedProduct.AffiliateURL != nil {
		product.AffiliateURL = *updatedProduct.AffiliateURL
	}

	if updatedProduct.CategoryID != nil {
		product.CategoryID = updatedProduct.CategoryID
	}

	updated, err := s.ProductRepository.UpdateProduct(product)

	if err != nil {
		return nil, errors.NewAppError(500, fmt.Sprintf("Failed to update product: %v", err.Error()))
	}

	return updated, nil
}

func (s *ProductServiceInstance) UpdateProductPhotos(product *models.Product) *errors.AppError {
	if _, err := s.ProductRepository.UpdateProduct(product); err != nil {
		return errors.NewAppError(500, fmt.Sprintf("Failed to update product photos: %v", err.Error()))
	}

	return nil
}

// Delete a product
// This function deletes a product from the repository by its ID
// It returns an error if any
func (s *ProductServiceInstance) DeleteProduct(productID string) *errors.AppError {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		return errors.NewAppError(404, fmt.Sprintf("Product not found: %v", err.Error()))
	}

	if product == nil {
		// return errors.New("product not found")
	}

	if err := s.ProductRepository.DeleteProduct(productID); err != nil {
		return errors.NewAppError(500, fmt.Sprintf("Failed to delete product: %v", err.Error()))
	}

	return nil
}

// Verify product ownership
// This function checks if the user has access to the product based on their merchant association
// It returns an error if the user does not have access or if the product is not found
func (s *ProductServiceInstance) VerifyProductOwnership(productID string, userID uint32) *errors.AppError {
	user, err := s.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		return errors.NewAppError(401, "unauthorized access")
	}

	if user.Role == "admin" {
		return nil
	}

	if len(user.Merchants) == 0 {
		return errors.NewAppError(403, "user does not have any merchants")
	}

	productMerchant, err := s.ProductRepository.FindMerchantByProductID(productID)
	if err != nil {
		return errors.NewAppError(404, fmt.Sprintf("Product not found: %v", err.Error()))
	}

	merchantID := productMerchant.ID
	for _, merchant := range user.Merchants {
		if merchant.ID == merchantID {
			return nil
		}
	}
	return errors.NewAppError(403, "user does not have any access to this product")
}
