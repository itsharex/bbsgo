package models

import "time"

type Follow struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"-"`
	FollowUserID uint      `gorm:"not null;index" json:"follow_user_id"`
	FollowUser   User      `gorm:"foreignKey:FollowUserID" json:"follow_user"`
	CreatedAt    time.Time `json:"created_at"`
}
