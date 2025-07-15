package services

import (
	"errors"
	"fmt"
	"senkou-catalyst-be/dtos"
	"senkou-catalyst-be/models"
	"senkou-catalyst-be/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, error)
	GetProductByID(productID string) (*models.Product, error)
	GetProductsByMerchantID(merchantID string) ([]*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, error)
	DeleteProduct(productID string) error
	VerifyProductOwnership(productID string, userID uint32) error
}

type productService struct {
	userReposiroty    repositories.UserRepository
	productRepository repositories.ProductRepository
}

func NewProductService(productRepository repositories.ProductRepository, userRepository repositories.UserRepository) ProductService {
	return &productService{
		productRepository: productRepository,
		userReposiroty:    userRepository,
	}
}

// Create a new product
// This function will be used to create a new product via repository
// It returns the created product and an error if any
func (s *productService) CreateProduct(product *dtos.CreateProductDTO, merchantID string) (*models.Product, error) {
	newProduct := &models.Product{
		ID:           uuid.New().String(),
		Title:        product.Title,
		Description:  product.Description,
		Price:        product.Price,
		AffiliateURL: product.AffiliateURL,
		CategoryID:   product.CategoryID,
		MerchantID:   merchantID,
	}

	storedProduct, err := s.productRepository.StoreProduct(newProduct)

	if err != nil {
		return nil, err
	}

	return storedProduct, nil
}

// Get a product by its ID
// This function retrieves a product from the repository by its ID
// It returns the product and an error if any
func (s *productService) GetProductByID(productID string) (*models.Product, error) {
	product, err := s.productRepository.FindProductByID(productID)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// Get products by merchant ID
// This function retrieves all products associated with a specific merchant ID
// It returns a slice of products and an error if any
func (s *productService) GetProductsByMerchantID(merchantID string) ([]*models.Product, error) {
	products, err := s.productRepository.FindProductsByMerchantID(merchantID)

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Get all products
// This function retrieves all products from the repository
// It returns a slice of products and an error if any
func (s *productService) GetAllProducts() ([]*models.Product, error) {
	products, err := s.productRepository.FindAllProducts()

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Update a product
// This function updates an existing product in the repository
// It takes an updated product request and a product ID as parameters
// It returns the updated product and an error if any
func (s *productService) UpdateProduct(updatedProduct *dtos.UpdateProductDTO, productID string) (*models.Product, error) {
	product, err := s.productRepository.FindProductByID(productID)

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
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

	updated, err := s.productRepository.UpdateProduct(product)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

// Delete a product
// This function deletes a product from the repository by its ID
// It returns an error if any
func (s *productService) DeleteProduct(productID string) error {
	product, err := s.productRepository.FindProductByID(productID)

	if err != nil {
		return err
	}

	if product == nil {
		return errors.New("product not found")
	}

	if err := s.productRepository.DeleteProduct(productID); err != nil {
		return err
	}

	return nil
}

func (s *productService) VerifyProductOwnership(productID string, userID uint32) error {
	user, err := s.userReposiroty.FindByID(userID)
	if err != nil || user == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "unauthorized access")
	}

	fmt.Println("FECTHED USER ID", user.ID)
	fmt.Println("USER_MERCHANTS", len(user.Merchants))

	if user.Role == "admin" {
		return nil
	}

	if len(user.Merchants) == 0 {
		return fiber.NewError(fiber.StatusForbidden, "user does not have any merchants")
	}

	productMerchant, err := s.productRepository.FindMerchantByProductID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	fmt.Println(productMerchant.Name)

	merchantID := productMerchant.ID
	for _, merchant := range user.Merchants {
		if merchant.ID == merchantID {
			return nil
		}
	}
	return fiber.NewError(fiber.StatusForbidden, "user does not have any access to this product")
}
