package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TopicID   uint           `gorm:"not null;index" json:"topic_id"`
	Topic     Topic          `gorm:"foreignKey:TopicID" json:"topic"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	ParentID  *uint          `gorm:"index" json:"parent_id"`
	Parent    *Post          `gorm:"foreignKey:ParentID" json:"-"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Children []Post `gorm:"foreignKey:ParentID" json:"children"`
	Likes    []Like `gorm:"foreignKey:TargetID" json:"-"`
}
