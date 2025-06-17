package vital

import (
	"ps_backend/model"

	"gorm.io/gorm"
)

// Repository handles database operations for vital signs.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository with the given database connection.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// CreateVital inserts a new VitalSign record.
func (r *Repository) CreateVital(entry *model.VitalSign) error {
	return r.db.Create(entry).Error
}

// GetVitalsByUser retrieves all VitalSign records for a specific user.
func (r *Repository) GetVitalsByUser(userID uint) ([]model.VitalSign, error) {
	var vitals []model.VitalSign
	if err := r.db.Where("user_id = ?", userID).Order("measured_at desc").Find(&vitals).Error; err != nil {
		return nil, err
	}
	return vitals, nil
}

// DeleteVital deletes a VitalSign record by its ID.
func (r *Repository) DeleteVital(id uint) error {
	return r.db.Delete(&model.VitalSign{}, id).Error
}
