package model

import "time"

// VitalSign represents a user's vital sign measurement record.
type VitalSign struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	HeartRate   int       `json:"heart_rate"`   // Beats per minute
	BreathRate  int       `json:"breath_rate"`  // Breaths per minute
	StressLevel int       `json:"stress_level"` // 0-100 scale
	MeasuredAt  time.Time `gorm:"not null" json:"measured_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
