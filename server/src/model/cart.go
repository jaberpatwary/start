package model

import (
	"time"
)

type Cart struct {
	ID        string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string     `gorm:"uniqueIndex;not null" json:"user_id"`
	Items     []CartItem `gorm:"constraint:OnDelete:CASCADE" json:"items"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        string   `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CartID    string   `gorm:"index;not null" json:"cart_id"`
	ProductID string   `gorm:"index;not null" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int      `gorm:"not null;default:1" json:"quantity"`
}
