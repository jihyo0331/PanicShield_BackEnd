package vital

import (
	"errors"
	"ps_backend/model"
	"time"

	"gorm.io/gorm"
)

// Service provides methods to manage vital signs.
type Service struct {
	repo *Repository
}

// NewService creates a new Service with the given gorm DB connection.
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

// CreateVital validates and creates a new vital sign entry.
func (s *Service) CreateVital(entry *model.VitalSign) error {
	if entry == nil {
		return errors.New("vital sign entry cannot be nil")
	}
	if entry.UserID == 0 {
		return errors.New("invalid user ID")
	}
	if entry.MeasuredAt.IsZero() {
		entry.MeasuredAt = time.Now()
	}
	return s.repo.CreateVital(entry)
}

// GetVitalsByUser retrieves all vital sign entries for a given user.
func (s *Service) GetVitalsByUser(userID uint) ([]model.VitalSign, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetVitalsByUser(userID)
}

// DeleteVital deletes a vital sign entry by its ID.
func (s *Service) DeleteVital(id uint) error {
	if id == 0 {
		return errors.New("invalid vital sign ID")
	}
	return s.repo.DeleteVital(id)
}
