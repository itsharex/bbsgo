package models

import "time"

type AntiSpamConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	Comment   string    `gorm:"size:255" json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserOperation struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Operation   string    `gorm:"size:50;not null;index" json:"operation"`
	TargetID    uint      `json:"target_id"`
	TargetType  string    `gorm:"size:50" json:"target_type"`
	ContentHash string    `gorm:"size:64" json:"content_hash"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

type ReputationLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Change    int       `gorm:"not null" json:"change"`
	Reason    string    `gorm:"size:255" json:"reason"`
	RelatedID uint      `json:"related_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ContentQuality struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	TargetID     uint      `gorm:"not null;index" json:"target_id"`
	TargetType   string    `gorm:"size:50;not null;index" json:"target_type"`
	QualityScore float64   `json:"quality_score"`
	IsLowQuality bool      `json:"is_low_quality"`
	Reasons      string    `gorm:"type:text" json:"reasons"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserBan struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Reason    string     `gorm:"size:255" json:"reason"`
	BanType   string     `gorm:"size:50" json:"ban_type"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
}
