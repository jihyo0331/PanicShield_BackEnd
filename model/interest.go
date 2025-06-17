package model

import "time"

type Interest struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}

type SubInterest struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	InterestID uint   `gorm:"not null;index" json:"interest_id"`
	Name       string `gorm:"size:255;not null" json:"name"`
}

// UserInterest represents a many-to-many mapping between users and interests.
type UserInterest struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	InterestID uint      `gorm:"not null;index" json:"interest_id"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserSubInterest represents a many-to-many mapping between users and sub-interests.
type UserSubInterest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	SubInterestID uint      `gorm:"not null;index" json:"sub_interest_id"`
	CreatedAt     time.Time `json:"created_at"`
}
