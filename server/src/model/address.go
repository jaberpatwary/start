package model

import (
	"time"
)

type Address struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID      string    `gorm:"index;not null" json:"user_id"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	Division    string    `json:"division"`
	District    string    `json:"district"`
	Thana       string    `json:"thana"`
	AddressLine string    `json:"address_line"`
	PostalCode  string    `json:"postal_code"`
	IsDefault   bool      `gorm:"default:false" json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
}
