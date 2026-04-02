package models

import "time"

type Tag struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Icon        string    `gorm:"size:10" json:"icon"`
	Description string    `gorm:"type:text" json:"description"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	UsageCount  int       `gorm:"default:0" json:"usage_count"`
	IsOfficial  bool      `gorm:"default:false" json:"is_official"`
	IsBanned    bool      `gorm:"default:false" json:"is_banned"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Tag) TableName() string {
	return "tags"
}
