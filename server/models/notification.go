package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Link      string    `gorm:"size:255" json:"link"`
	IsRead    bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
