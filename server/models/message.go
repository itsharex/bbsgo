package models

import "time"

type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	FromUserID uint      `gorm:"not null;index" json:"from_user_id"`
	FromUser   User      `gorm:"foreignKey:FromUserID" json:"from_user"`
	ToUserID   uint      `gorm:"not null;index" json:"to_user_id"`
	ToUser     User      `gorm:"foreignKey:ToUserID" json:"to_user"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	IsRead     bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}
