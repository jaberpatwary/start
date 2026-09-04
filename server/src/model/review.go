package model

import (
	"time"
)

const (
	ReviewPending  = "PENDING"
	ReviewApproved = "APPROVED"
	ReviewRejected = "REJECTED"
)

type Review struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProductID string    `gorm:"index;not null" json:"product_id"`
	UserID    string    `gorm:"index;not null" json:"user_id"`
	User      *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Rating    int       `gorm:"not null" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	Status    string    `gorm:"default:PENDING;index" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
