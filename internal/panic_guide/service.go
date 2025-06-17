package panic_guide

import (
	"errors"
	"ps_backend/model"

	"gorm.io/gorm"
)

// Repository provides database access for panic guides.
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new panic guide Repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetAllPanicGuides retrieves all panic guides.
func (r *Repository) GetAllPanicGuides() ([]model.PanicGuide, error) {
	var guides []model.PanicGuide
	if err := r.db.Find(&guides).Error; err != nil {
		return nil, err
	}
	return guides, nil
}

// CreatePanicGuide inserts a new panic guide with given title and description.
func (r *Repository) CreatePanicGuide(title, description string) (*model.PanicGuide, error) {
	guide := &model.PanicGuide{Title: title, Description: description}
	if err := r.db.Create(guide).Error; err != nil {
		return nil, err
	}
	return guide, nil
}

// BookmarkPanicGuide creates a user_panic_guides entry.
func (r *Repository) BookmarkPanicGuide(userID, guideID uint) error {
	entry := &model.UserPanicGuide{UserID: userID, PanicGuideID: guideID}
	return r.db.Create(entry).Error
}

// GetUserBookmarkedGuides returns guides bookmarked by a user.
func (r *Repository) GetUserBookmarkedGuides(userID uint) ([]model.PanicGuide, error) {
	var guides []model.PanicGuide
	if err := r.db.
		Table("panic_guides").
		Select("panic_guides.*").
		Joins("join user_panic_guides on user_panic_guides.panic_guide_id = panic_guides.id").
		Where("user_panic_guides.user_id = ?", userID).
		Scan(&guides).Error; err != nil {
		return nil, err
	}
	return guides, nil
}

// Service provides business logic for panic guides.
type Service struct {
	repo *Repository
}

// NewService creates a new panic guide Service using the given DB connection.
func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

// GetAll retrieves all panic guides.
func (s *Service) GetAll() ([]model.PanicGuide, error) {
	return s.repo.GetAllPanicGuides()
}

// Create adds a new panic guide with given title and description.
func (s *Service) Create(title, description string) (*model.PanicGuide, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if description == "" {
		return nil, errors.New("description cannot be empty")
	}
	return s.repo.CreatePanicGuide(title, description)
}

// Bookmark marks a panic guide as bookmarked for a user.
func (s *Service) Bookmark(userID, guideID uint) error {
	if userID == 0 || guideID == 0 {
		return errors.New("userID and guideID must be provided")
	}
	return s.repo.BookmarkPanicGuide(userID, guideID)
}

// GetBookmarks retrieves all bookmarked guides for a user.
func (s *Service) GetBookmarks(userID uint) ([]model.PanicGuide, error) {
	if userID == 0 {
		return nil, errors.New("userID must be provided")
	}
	return s.repo.GetUserBookmarkedGuides(userID)
}
