package model

import "time"

// PanicGuide represents a method to cope with panic attacks.
type PanicGuide struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:128;not null"`
	Description string `gorm:"type:text;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserPanicGuide represents a user's bookmarked panic guide.
type UserPanicGuide struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	PanicGuideID uint      `gorm:"not null;index"`
	BookmarkedAt time.Time `gorm:"autoCreateTime"`
}
