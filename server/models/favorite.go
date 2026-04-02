package models

import "time"

type Favorite struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	TopicID   uint      `gorm:"not null;index" json:"topic_id"`
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"topic"`
	CreatedAt time.Time `json:"created_at"`
}
