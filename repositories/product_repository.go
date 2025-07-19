package repositories

import (
	"errors"
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	StoreProduct(product *models.Product) (*models.Product, error)
	FindProductByID(productID string) (*models.Product, error)
	FindProductsByMerchantID(merchantID string) ([]*models.Product, error)
	FindAllProducts() ([]*models.Product, error)
	FindMerchantByProductID(productID string) (*models.Merchant, error)
	UpdateProduct(updatedProduct *models.Product) (*models.Product, error)
	DeleteProduct(productID string) error
}

type ProductRepositoryInstance struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryInstance{
		DB: db,
	}
}

// Store a new product
// This function will be used to store a new product in the database
// It takes a product model and a merchant ID as parameters
// It returns the stored product and an error if any
func (r *ProductRepositoryInstance) StoreProduct(product *models.Product) (*models.Product, error) {
	if err := r.DB.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// Get a product by its ID
// This function retrieves a product from the database by its ID
// It takes a product ID as a parameter
// It returns the product and an error if any
func (r *ProductRepositoryInstance) FindProductByID(productID string) (*models.Product, error) {
	product := new(models.Product)

	if err := r.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// Get products by merchant ID
// This function retrieves all products associated with a specific merchant ID
// It takes a merchant ID as a parameter
// It returns a slice of products and an error if any
func (r *ProductRepositoryInstance) FindProductsByMerchantID(merchantID string) ([]*models.Product, error) {
	products := make([]*models.Product, 0)

	if err := r.DB.Where("merchant_id = ?", merchantID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// Get all products by administrator only
// This function retrieves all products from the database
// It returns a slice of products and an error if any
func (r *ProductRepositoryInstance) FindAllProducts() ([]*models.Product, error) {
	products := make([]*models.Product, 0)

	if err := r.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// Get the merchant associated with a product by its ID
// This function retrieves the merchant associated with a specific product ID
// It returns the merchant and an error if any
func (r *ProductRepositoryInstance) FindMerchantByProductID(productID string) (*models.Merchant, error) {
	merchant := new(models.Merchant)

	if err := r.DB.Table("merchants").Select("merchants.*").
		Joins("JOIN products ON products.merchant_id = merchants.id").
		Where("products.id = ?", productID).First(&merchant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return merchant, nil
}

// Update a product
// This function updates an existing product in the database
// It takes a product model as a parameter
// It returns the updated product and an error if any
func (r *ProductRepositoryInstance) UpdateProduct(updatedProduct *models.Product) (*models.Product, error) {
	if err := r.DB.Save(updatedProduct).Error; err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

// Delete a product by its ID
// This function deletes a product from the database by its ID
// It takes a product ID as a parameter
// It returns an error if any
func (r *ProductRepositoryInstance) DeleteProduct(productID string) error {
	if err := r.DB.Where("id = ?", productID).Delete(&models.Product{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}

		return err
	}

	return nil
}
