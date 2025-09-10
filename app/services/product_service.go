package services

import (
	"fmt"
	"senkou-catalyst-be/app/dtos"
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/platform/errors"
	"senkou-catalyst-be/repositories"
	"senkou-catalyst-be/utils/query"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, *errors.CustomError)
	GetProductByID(productID string) (*models.Product, *errors.CustomError)
	GetProductsByMerchantID(merchantID string) ([]*models.Product, *errors.CustomError)
	GetAllProducts(params *query.QueryParams) ([]*models.Product, *query.PaginationResponse, *errors.CustomError)
	GetProductsByMerchantUsername(username string) ([]*models.Product, *errors.CustomError)
	UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, *errors.CustomError)
	UpdateProductPhotos(product *models.Product) *errors.CustomError
	DeleteProduct(productID string) *errors.CustomError
	VerifyProductOwnership(productID string, userID uint32) *errors.CustomError
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
func (s *ProductServiceInstance) CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, *errors.CustomError) {
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
		return nil, errors.Internal("Failed to create product", err.Error())
	}

	return storedProduct, nil
}

// Get a product by its ID
// This function retrieves a product from the repository by its ID
// It returns the product and an error if any
func (s *ProductServiceInstance) GetProductByID(productID string) (*models.Product, *errors.CustomError) {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("Product not found")
		}

		return nil, errors.Internal("Failed to retrieve product", err.Error())
	}

	return product, nil
}

// Get products by merchant username
// This function retrieves all products associated with a specific merchant username
// It returns a slice of products and an error if any
func (s *ProductServiceInstance) GetProductsByMerchantUsername(username string) ([]*models.Product, *errors.CustomError) {
	products, err := s.ProductRepository.FindProductsByMerchantUsername(username)
	if err != nil {
		return nil, errors.NotFound(fmt.Sprintf("Products not found for merchant username %s", username))
	}

	return products, nil
}

// Get products by merchant ID
// This function retrieves all products associated with a specific merchant ID
// It returns a slice of products and an error if any
func (s *ProductServiceInstance) GetProductsByMerchantID(merchantID string) ([]*models.Product, *errors.CustomError) {
	products, err := s.ProductRepository.FindProductsByMerchantID(merchantID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("No products found for this merchant")
		}

		return nil, errors.Internal("Failed to retrieve products", err.Error())
	}

	return products, nil
}

// Get all products
// This function retrieves all products from the repository
// It returns a slice of products and an error if any
func (s *ProductServiceInstance) GetAllProducts(params *query.QueryParams) ([]*models.Product, *query.PaginationResponse, *errors.CustomError) {
	products, total, err := s.ProductRepository.FindAllProducts(params)

	if err != nil {
		return nil, nil, errors.Internal("Failed to retrieve products", err.Error())
	}

	pagination := query.CalculatePagination(params.Page, params.Limit, total)

	return products, pagination, nil
}

// Update a product
// This function updates an existing product in the repository
// It takes an updated product request and a product ID as parameters
// It returns the updated product and an error if any
func (s *ProductServiceInstance) UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, *errors.CustomError) {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		return nil, errors.NotFound("Product not found")
	}

	if product == nil {
		return nil, errors.NotFound("Product not found")
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
		return nil, errors.Internal("Failed to update product", err.Error())
	}

	return updated, nil
}

func (s *ProductServiceInstance) UpdateProductPhotos(product *models.Product) *errors.CustomError {
	if _, err := s.ProductRepository.UpdateProduct(product); err != nil {
		return errors.Internal("Failed to update product photos", err.Error())
	}

	return nil
}

// Delete a product
// This function deletes a product from the repository by its ID
// It returns an error if any
func (s *ProductServiceInstance) DeleteProduct(productID string) *errors.CustomError {
	product, err := s.ProductRepository.FindProductByID(productID)

	if err != nil {
		return errors.NotFound("Product not found")
	}

	if product == nil {
		return errors.NotFound("Product not found")
	}

	if err := s.ProductRepository.DeleteProduct(productID); err != nil {
		return errors.Internal("Failed to delete product", err.Error())
	}

	return nil
}

// Verify product ownership
// This function checks if the user has access to the product based on their merchant association
// It returns an error if the user does not have access or if the product is not found
func (s *ProductServiceInstance) VerifyProductOwnership(productID string, userID uint32) *errors.CustomError {
	user, err := s.UserRepository.FindByID(userID)
	if err != nil || user == nil {
		return errors.Unauthorized("Unauthorized access")
	}

	if user.Role == "admin" {
		return nil
	}

	if len(user.Merchants) == 0 {
		return errors.Forbidden("User does not have any merchants")
	}

	productMerchant, err := s.ProductRepository.FindMerchantByProductID(productID)
	if err != nil {
		return errors.NotFound("Product not found")
	}

	merchantID := productMerchant.ID
	for _, merchant := range user.Merchants {
		if merchant.ID == merchantID {
			return nil
		}
	}
	return errors.Forbidden("User does not have access to this product")
}
