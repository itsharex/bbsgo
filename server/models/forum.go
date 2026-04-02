package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type IntSlice []int

func (s IntSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *IntSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type Forum struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	SortOrder    int       `gorm:"default:0" json:"sort_order"`
	Icon         string    `gorm:"size:255" json:"icon"`
	ModeratorIDs IntSlice  `gorm:"type:json" json:"moderator_ids"`
	AllowPost    bool      `gorm:"default:true" json:"allow_post"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Topics []Topic `gorm:"foreignKey:ForumID" json:"-"`
}
