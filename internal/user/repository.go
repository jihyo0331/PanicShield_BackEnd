package user

import (
	"ps_backend/model"

	"gorm.io/gorm"
)

// Repository handles CRUD operations for User entities.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetByID retrieves a user by its ID.
func (r *Repository) GetByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username.
func (r *Repository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername checks if a username already exists.
func (r *Repository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByPhone checks if a phone number already exists.
func (r *Repository) ExistsByPhone(phone string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("phone_number = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create adds a new user to the database.
func (r *Repository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update saves the updated user data.
func (r *Repository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user by ID.
func (r *Repository) Delete(userID uint) error {
	return r.db.Delete(&model.User{}, userID).Error
}
