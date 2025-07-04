package repositories

import (
	"senkou-catalyst-be/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	Create(user models.User) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Create(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}
