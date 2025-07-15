package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindAll() (*[]models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(userID uint32) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Find all users in the database
// Returns a slice of User models or an error if the operation fails
func (r *userRepository) FindAll() (*[]models.User, error) {
	var users *[]models.User
	err := r.db.Find(&users).Error
	return users, err
}

// Create a new user in the database
// Returns the created user or an error if any
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

// Find a user by its email
// Returns the user model or an error if not found
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Find a user by its ID
// Returns true if the user exists, false if not found, or an error if any
func (r *userRepository) FindByID(userID uint32) (*models.User, error) {
	user := new(models.User)

	if err := r.db.Preload("Merchants").First(user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}

		return nil, err
	}

	return user, nil
}
