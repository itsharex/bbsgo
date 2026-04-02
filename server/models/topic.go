package models

import (
	"time"

	"gorm.io/gorm"
)

type Topic struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user"`
	ForumID      uint           `gorm:"not null;index" json:"forum_id"`
	Forum        Forum          `gorm:"foreignKey:ForumID" json:"forum"`
	IsPinned     bool           `gorm:"default:false;index" json:"is_pinned"`
	IsLocked     bool           `gorm:"default:false" json:"is_locked"`
	IsEssence    bool           `gorm:"default:false" json:"is_essence"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	ViewCount    int            `gorm:"default:0" json:"view_count"`
	ReplyCount   int            `gorm:"default:0" json:"reply_count"`
	LastReplyAt  *time.Time     `json:"last_reply_at"`
	AllowComment bool           `gorm:"default:true" json:"allow_comment"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Posts     []Post     `gorm:"foreignKey:TopicID" json:"-"`
	Likes     []Like     `gorm:"foreignKey:TargetID" json:"-"`
	Favorites []Favorite `gorm:"foreignKey:TopicID" json:"-"`
	Tags      []Tag      `gorm:"many2many:topic_tags;" json:"tags"`
}
