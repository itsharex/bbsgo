package models

import "time"

type UserBadge struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	BadgeID   uint      `gorm:"not null;index" json:"badge_id"`
	AwardedAt time.Time `json:"awarded_at"`
}

func (UserBadge) TableName() string {
	return "user_badges"
}
