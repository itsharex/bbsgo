package models

import "time"

type Report struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ReporterID  uint      `gorm:"not null;index" json:"reporter_id"`
	Reporter    User      `gorm:"foreignKey:ReporterID" json:"reporter"`
	TargetType  string    `gorm:"size:20;not null" json:"target_type"`
	TargetID    uint      `gorm:"not null" json:"target_id"`
	Reason      string    `gorm:"type:text;not null" json:"reason"`
	Status      int       `gorm:"default:0;index" json:"status"`
	HandledAt   *time.Time `json:"handled_at"`
	HandlerID   *uint     `json:"handler_id"`
	Handler     *User     `gorm:"foreignKey:HandlerID" json:"handler"`
	CreatedAt   time.Time `json:"created_at"`
}
