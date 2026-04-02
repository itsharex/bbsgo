package models

import "time"

type ForumCategory struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Icon        string    `gorm:"size:10" json:"icon"`
	Description string    `gorm:"size:200" json:"description"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
