package model

import (
	"time"
)

type Banner struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Image     string    `json:"image"`
	Link      string    `json:"link"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
}
