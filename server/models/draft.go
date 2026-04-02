package models

import "time"

type Draft struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Title     string    `gorm:"size:200" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	ForumID   uint      `json:"forum_id"`
	Tags      IntSlice  `gorm:"type:json" json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
