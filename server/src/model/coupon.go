package model

import (
	"time"
)

const (
	CouponPercent = "PERCENT"
	CouponFixed   = "FIXED"
)

type Coupon struct {
	ID        string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Code      string     `gorm:"uniqueIndex;not null" json:"code"`
	Type      string     `gorm:"not null" json:"type"`
	Value     int        `json:"value"`
	MinOrder  int        `json:"min_order"`
	MaxUses   int        `json:"max_uses"`
	UsedCount int        `gorm:"default:0" json:"used_count"`
	ExpiresAt *time.Time `json:"expires_at"`
	Active    bool       `gorm:"default:true" json:"active"`
	CreatedAt time.Time  `json:"created_at"`
}
