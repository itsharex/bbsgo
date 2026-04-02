package models

import "time"

type Badge struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Icon        string    `gorm:"size:255" json:"icon"`
	Condition   string    `gorm:"type:text" json:"condition"`
	CreatedAt   time.Time `json:"created_at"`
}
