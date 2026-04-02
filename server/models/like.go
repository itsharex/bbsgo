package models

import "time"

type Like struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	TargetType string    `gorm:"size:20;not null;index" json:"target_type"`
	TargetID   uint      `gorm:"not null;index" json:"target_id"`
	CreatedAt  time.Time `json:"created_at"`
}
