package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         int            `gorm:"default:0" json:"role"`
	Avatar       string         `gorm:"size:255" json:"avatar"`
	Background   string         `gorm:"size:255" json:"background"`
	Signature    string         `gorm:"size:255" json:"signature"`
	Intro        string         `gorm:"type:text" json:"intro"`
	Credits      int            `gorm:"default:0" json:"credits"`
	Level        int            `gorm:"default:1" json:"level"`
	LastSignAt   *time.Time     `json:"last_sign_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Topics         []Topic        `gorm:"foreignKey:UserID" json:"-"`
	Posts          []Post         `gorm:"foreignKey:UserID" json:"-"`
	Likes          []Like         `gorm:"foreignKey:UserID" json:"-"`
	Favorites      []Favorite     `gorm:"foreignKey:UserID" json:"-"`
	Follows        []Follow       `gorm:"foreignKey:UserID" json:"-"`
	Followers      []Follow       `gorm:"foreignKey:FollowUserID" json:"-"`
	SentMessages   []Message      `gorm:"foreignKey:FromUserID" json:"-"`
	ReceivedMessages []Message    `gorm:"foreignKey:ToUserID" json:"-"`
	Notifications  []Notification `gorm:"foreignKey:UserID" json:"-"`
	Reports        []Report       `gorm:"foreignKey:ReporterID" json:"-"`
	Drafts         []Draft        `gorm:"foreignKey:UserID" json:"-"`
	UserBadges     []UserBadge    `gorm:"foreignKey:UserID" json:"-"`
	FollowedTopics []Topic        `gorm:"many2many:user_follow_topics" json:"-"`
}
