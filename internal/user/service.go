package user

import (
	"errors"
	"ps_backend/model"

	"gorm.io/gorm"
)

// Service provides user-related business logic.
type Service struct {
	repo *Repository
}

// NewService constructs a new Service with a GORM DB.
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

// GetByID returns the user with the specified ID.
func (s *Service) GetByID(userID uint) (*model.User, error) {
	if userID == 0 {
		return nil, errors.New("userID must be provided")
	}
	return s.repo.GetByID(userID)
}

// GetByUsername returns the user matching the given username.
func (s *Service) GetByUsername(username string) (*model.User, error) {
	if username == "" {
		return nil, errors.New("username must be provided")
	}
	return s.repo.GetByUsername(username)
}

// Create creates a new user record.
// Expects the user.PasswordHash already set.
func (s *Service) Create(user *model.User) error {
	if user == nil {
		return errors.New("user must be provided")
	}
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.PasswordHash == "" {
		return errors.New("password hash cannot be empty")
	}
	if user.PhoneNumber == "" {
		return errors.New("phone number cannot be empty")
	}
	return s.repo.Create(user)
}

// Update updates existing user data.
func (s *Service) Update(user *model.User) error {
	if user == nil || user.ID == 0 {
		return errors.New("user and userID must be provided")
	}
	return s.repo.Update(user)
}

// Delete removes a user by ID.
func (s *Service) Delete(userID uint) error {
	if userID == 0 {
		return errors.New("userID must be provided")
	}
	return s.repo.Delete(userID)
}
