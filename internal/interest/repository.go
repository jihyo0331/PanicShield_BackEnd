package interest

import (
	"ps_backend/model"

	"gorm.io/gorm"
)

// Repository provides access to interest and sub-interest storage.
type Repository struct {
	db *gorm.DB
}

// NewRepository returns a new instance of Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetAllInterests retrieves all interests from the database.
func (r *Repository) GetAllInterests() ([]model.Interest, error) {
	var interests []model.Interest
	if err := r.db.Find(&interests).Error; err != nil {
		return nil, err
	}
	return interests, nil
}

// GetSubInterests retrieves all sub-interests associated with a given interest ID.
func (r *Repository) GetSubInterests(interestID uint) ([]model.SubInterest, error) {
	var subInterests []model.SubInterest
	if err := r.db.Where("interest_id = ?", interestID).Find(&subInterests).Error; err != nil {
		return nil, err
	}
	return subInterests, nil
}

// CreateInterest creates a new interest with the given name.
func (r *Repository) CreateInterest(name string) (*model.Interest, error) {
	interest := &model.Interest{Name: name}
	if err := r.db.Create(interest).Error; err != nil {
		return nil, err
	}
	return interest, nil
}

// CreateSubInterest creates a new sub-interest under the specified interest.
func (r *Repository) CreateSubInterest(interestID uint, name string) (*model.SubInterest, error) {
	subInterest := &model.SubInterest{InterestID: interestID, Name: name}
	if err := r.db.Create(subInterest).Error; err != nil {
		return nil, err
	}
	return subInterest, nil
}

// AssignInterestToUser assigns an interest to a user using a join table.
func (r *Repository) AssignInterestToUser(userID, interestID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Check if already assigned
		var count int64
		err := tx.Table("user_interests").
			Where("user_id = ? AND interest_id = ?", userID, interestID).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return nil // Already assigned, no-op
		}
		// Insert assignment
		return tx.Table("user_interests").Create(map[string]interface{}{
			"user_id":     userID,
			"interest_id": interestID,
		}).Error
	})
}

// RemoveInterestFromUser removes an interest assignment from a user.
func (r *Repository) RemoveInterestFromUser(userID, interestID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Table("user_interests").
			Where("user_id = ? AND interest_id = ?", userID, interestID).
			Delete(nil).Error
	})
}

// GetUserInterests returns all interests assigned to the user.
func (r *Repository) GetUserInterests(userID uint) ([]model.Interest, error) {
	var interests []model.Interest
	err := r.db.
		Joins("JOIN user_interests ON user_interests.interest_id = interests.id").
		Where("user_interests.user_id = ?", userID).
		Find(&interests).Error
	if err != nil {
		return nil, err
	}
	return interests, nil
}

// GetUserSubInterests returns all sub-interests assigned to the user.
func (r *Repository) GetUserSubInterests(userID uint) ([]model.SubInterest, error) {
	var subInterests []model.SubInterest
	err := r.db.
		Joins("JOIN user_sub_interests ON user_sub_interests.sub_interest_id = sub_interests.id").
		Where("user_sub_interests.user_id = ?", userID).
		Find(&subInterests).Error
	if err != nil {
		return nil, err
	}
	return subInterests, nil
}
