package model

import (
	"time"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

type User struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name         string     `gorm:"not null" json:"name"`
	Email        string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone        string     `json:"phone"`
	PasswordHash string     `json:"-"`
	Role         string     `gorm:"default:USER" json:"role"`
	Avatar       string     `json:"avatar"`
	IsBlocked    bool       `gorm:"default:false" json:"is_blocked"`
	Addresses    []Address  `gorm:"constraint:OnDelete:CASCADE" json:"addresses,omitempty"`
	Orders       []Order    `json:"orders,omitempty"`
	Reviews      []Review   `json:"reviews,omitempty"`
	Wishlist     []Wishlist `gorm:"constraint:OnDelete:CASCADE" json:"wishlist,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
