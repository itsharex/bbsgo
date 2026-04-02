package models

import "time"

type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Code      string    `gorm:"not null" json:"code"`
	Type      string    `gorm:"not null;default:'register'" json:"type"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (VerificationCode) TableName() string {
	return "verification_codes"
}
