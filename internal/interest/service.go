package interest

import (
	"errors"
	"ps_backend/model"

	"gorm.io/gorm"
)

// Service provides business logic for interests and sub-interests.
type Service struct {
	repo *Repository
}

// NewService creates a new interest Service with given DB.
func NewService(db *gorm.DB) *Service {
	return &Service{
		repo: NewRepository(db),
	}
}

// GetAll returns all interests.
func (s *Service) GetAll() ([]model.Interest, error) {
	return s.repo.GetAllInterests()
}

// GetSub returns all sub-interests for a given interest ID.
func (s *Service) GetSub(interestID uint) ([]model.SubInterest, error) {
	if interestID == 0 {
		return nil, errors.New("interestID must be provided")
	}
	return s.repo.GetSubInterests(interestID)
}

// Create adds a new interest with the given name.
func (s *Service) Create(name string) (*model.Interest, error) {
	if name == "" {
		return nil, errors.New("interest name cannot be empty")
	}
	return s.repo.CreateInterest(name)
}

// CreateSub adds a new sub-interest under the given interest.
func (s *Service) CreateSub(interestID uint, name string) (*model.SubInterest, error) {
	if interestID == 0 {
		return nil, errors.New("interestID must be provided")
	}
	if name == "" {
		return nil, errors.New("sub-interest name cannot be empty")
	}
	// Verify interest exists
	all, err := s.repo.GetAllInterests()
	if err != nil {
		return nil, err
	}
	found := false
	for _, inter := range all {
		if inter.ID == interestID {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("parent interest not found")
	}
	return s.repo.CreateSubInterest(interestID, name)
}

// AssignToUser links an interest to a user.
func (s *Service) AssignToUser(userID, interestID uint) error {
	if userID == 0 || interestID == 0 {
		return errors.New("userID and interestID must be provided")
	}
	return s.repo.AssignInterestToUser(userID, interestID)
}

// RemoveFromUser unlinks an interest from a user.
func (s *Service) RemoveFromUser(userID, interestID uint) error {
	if userID == 0 || interestID == 0 {
		return errors.New("userID and interestID must be provided")
	}
	return s.repo.RemoveInterestFromUser(userID, interestID)
}

// GetUserInterests returns interests assigned to a user.
func (s *Service) GetUserInterests(userID uint) ([]model.Interest, error) {
	if userID == 0 {
		return nil, errors.New("userID must be provided")
	}
	return s.repo.GetUserInterests(userID)
}

// GetUserSubInterests returns sub-interests assigned to a user.
func (s *Service) GetUserSubInterests(userID uint) ([]model.SubInterest, error) {
	if userID == 0 {
		return nil, errors.New("userID must be provided")
	}
	return s.repo.GetUserSubInterests(userID)
}
