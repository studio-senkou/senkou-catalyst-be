package repositories

import (
	"senkou-catalyst-be/app/models"
	"senkou-catalyst-be/utils/query"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	FindAll(params *query.QueryParams) (*[]models.User, int64, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(userID uint32) (*models.User, error)
	Update(user *models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Find all users in the database
// Returns a slice of User models or an error if the operation fails
func (r *userRepository) FindAll(params *query.QueryParams) (*[]models.User, int64, error) {
	users := make([]models.User, 0)
	var totalRecords int64

	queryBuilder := query.NewQueryBuilder(r.db.Model(&models.User{})).SetAllowedSorts(map[string]string{
		"name":       "name",
		"email":      "email",
		"created_at": "created_at",
	}).SetSearchFields([]string{"name", "email"})

	baseQuery := queryBuilder.BuildQuery(params)

	if err := baseQuery.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	paginatedQuery := queryBuilder.ApplyPagination(baseQuery, params)
	if err := paginatedQuery.Preload("Merchants").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return &users, totalRecords, nil
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
	user := new(models.User)

	err := r.db.Where("email = ?", email).First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
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

// Update a user in the database
// Returns the updated user or an error if any
func (r *userRepository) Update(user *models.User) (*models.User, error) {
	if err := r.db.Model(user).Updates(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
