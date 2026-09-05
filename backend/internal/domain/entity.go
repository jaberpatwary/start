package domain

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"

	ProductActive     = "ACTIVE"
	ProductDraft      = "DRAFT"
	ProductOutOfStock = "OUT_OF_STOCK"

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

	ReviewPending  = "PENDING"
	ReviewApproved = "APPROVED"
	ReviewRejected = "REJECTED"

	CouponPercent = "PERCENT"
	CouponFixed   = "FIXED"
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

type Category struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	Icon      string    `json:"icon"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type Brand struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID               string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name             string         `gorm:"not null;index" json:"name"`
	Slug             string         `gorm:"uniqueIndex;not null" json:"slug"`
	SKU              string         `gorm:"uniqueIndex;not null" json:"sku"`
	CategoryID       string         `gorm:"index;not null" json:"category_id"`
	Category         *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	BrandID          string         `gorm:"index;not null" json:"brand_id"`
	Brand            *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Price            int            `gorm:"not null" json:"price"`
	DiscountPrice    *int           `json:"discount_price"`
	Stock            int            `gorm:"default:0" json:"stock"`
	ShortDescription string         `json:"short_description"`
	Description      string         `gorm:"type:text" json:"description"`
	Images           []ProductImage `gorm:"constraint:OnDelete:CASCADE" json:"images,omitempty"`
	Specs            datatypes.JSON `json:"specs"`
	Tags             pq.StringArray `gorm:"type:text[]" json:"tags"`
	Status           string         `gorm:"default:ACTIVE;index" json:"status"`
	Featured         bool           `gorm:"default:false" json:"featured"`
	Sold             int            `gorm:"default:0" json:"sold"`
	RatingAvg        float64        `gorm:"default:0" json:"rating_avg"`
	RatingCount      int            `gorm:"default:0" json:"rating_count"`
	Reviews          []Review       `json:"reviews,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ProductID string `gorm:"index;not null" json:"product_id"`
	URL       string `gorm:"not null" json:"url"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

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

type Wishlist struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"index;not null" json:"user_id"`
	ProductID string    `gorm:"index;not null" json:"product_id"`
	Product   *Product  `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

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

type Setting struct {
	Key   string         `gorm:"primaryKey" json:"key"`
	Value datatypes.JSON `json:"value"`
}
