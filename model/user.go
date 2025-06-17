package model

import "time"

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Username      string `gorm:"unique;not null"`
	PasswordHash  string `gorm:"not null"`
	PhoneNumber   string `gorm:"unique;not null"`
	Verified      bool   `gorm:"default:false"`
	SpeakingStyle string // 예: 반말, 존댓말
	Tone          string // 예: 유머, 딱딱함 등
	CreatedAt     time.Time
}

type ChatbotLog struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Message   string // 사용자의 질문 또는 챗봇의 답변
	Sender    string // 'user' 또는 'bot'
	CreatedAt time.Time
}
