package model

import (
	"time"

	"gorm.io/datatypes"
)

const (
	PaymentCOD      = "COD"
	PaymentBKash    = "BKASH"
	PaymentNagad    = "NAGAD"
	PaymentRocket   = "ROCKET"
	PaymentCard     = "CARD"
	PaymentBank     = "BANK"
	PaymentUnpaid   = "UNPAID"
	PaymentPaid     = "PAID"
	PaymentRefunded = "REFUNDED"

	OrderPending    = "PENDING"
	OrderConfirmed  = "CONFIRMED"
	OrderProcessing = "PROCESSING"
	OrderShipped    = "SHIPPED"
	OrderDelivered  = "DELIVERED"
	OrderCancelled  = "CANCELLED"
)

type Order struct {
	ID               string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	OrderNumber      string         `gorm:"uniqueIndex;not null" json:"order_number"`
	UserID           string         `gorm:"index;not null" json:"user_id"`
	User             *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items            []OrderItem    `gorm:"constraint:OnDelete:CASCADE" json:"items"`
	ShippingName     string         `json:"shipping_name"`
	ShippingPhone    string         `json:"shipping_phone"`
	ShippingDivision string         `json:"shipping_division"`
	ShippingDistrict string         `json:"shipping_district"`
	ShippingThana    string         `json:"shipping_thana"`
	ShippingAddress  string         `json:"shipping_address"`
	ShippingPostal   string         `json:"shipping_postal"`
	PaymentMethod    string         `gorm:"not null" json:"payment_method"`
	PaymentStatus    string         `gorm:"default:UNPAID" json:"payment_status"`
	Subtotal         int            `json:"subtotal"`
	Discount         int            `json:"discount"`
	ShippingFee      int            `json:"shipping_fee"`
	Total            int            `json:"total"`
	Status           string         `gorm:"default:PENDING;index" json:"status"`
	TrackingNumber   string         `json:"tracking_number"`
	CouponCode       string         `json:"coupon_code"`
	Note             string         `json:"note"`
	StatusHistory    datatypes.JSON `json:"status_history"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type OrderItem struct {
	ID        string   `gorm:"primaryKey;type:varchar(36)" json:"id"`
	OrderID   string   `gorm:"index;not null" json:"order_id"`
	ProductID string   `gorm:"index;not null" json:"product_id"`
	Product   *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Name      string   `json:"name"`
	Price     int      `json:"price"`
	Quantity  int      `json:"quantity"`
	Image     string   `json:"image"`
}
