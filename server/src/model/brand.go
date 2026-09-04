package model

import (
	"time"
)

type Brand struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
}
