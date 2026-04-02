package models

import "time"

type Announcement struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsPinned  bool      `gorm:"default:false" json:"is_pinned"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
