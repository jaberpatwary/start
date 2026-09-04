package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/datatypes"
)

const (
	ProductActive     = "ACTIVE"
	ProductDraft      = "DRAFT"
	ProductOutOfStock = "OUT_OF_STOCK"
)

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
