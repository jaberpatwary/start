package domain

import (
	"context"
)

type ProductFilter struct {
	Category string
	Brand    string
	Search   string
	MinPrice int
	MaxPrice int
	Rating   float64
	Sort     string
	Page     int
	Limit    int
	Featured *bool
	Status   string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	GetAll(ctx context.Context, page, limit int) ([]User, int64, error)
	ToggleBlock(ctx context.Context, id string) (*User, error)
	AddAddress(ctx context.Context, addr *Address) error
	DeleteAddress(ctx context.Context, id, userID string) error
}

type ProductRepository interface {
	GetAll(ctx context.Context, filter ProductFilter) ([]Product, int64, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	Create(ctx context.Context, product *Product, images []string) error
	Update(ctx context.Context, product *Product, images []string) error
	Delete(ctx context.Context, id string) error
	GetReviews(ctx context.Context, productID string) ([]Review, error)
	CreateReview(ctx context.Context, review *Review) error
}

type CatalogRepository interface {
	GetCategories(ctx context.Context) ([]Category, error)
	GetCategoryByID(ctx context.Context, id string) (*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id string) error

	GetBrands(ctx context.Context) ([]Brand, error)
	GetBrandByID(ctx context.Context, id string) (*Brand, error)
	CreateBrand(ctx context.Context, brand *Brand) error
	UpdateBrand(ctx context.Context, brand *Brand) error
	DeleteBrand(ctx context.Context, id string) error

	GetBanners(ctx context.Context) ([]Banner, error)
	CreateBanner(ctx context.Context, banner *Banner) error
	DeleteBanner(ctx context.Context, id string) error

	GetCouponByCode(ctx context.Context, code string) (*Coupon, error)
	CreateCoupon(ctx context.Context, coupon *Coupon) error
	DeleteCoupon(ctx context.Context, id string) error
	IncrementCouponUsage(ctx context.Context, code string) error
}

type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID string) (*Cart, error)
	AddItem(ctx context.Context, userID, productID string, quantity int) (*Cart, error)
	UpdateItem(ctx context.Context, userID, itemID string, quantity int) (*Cart, error)
	RemoveItem(ctx context.Context, userID, itemID string) (*Cart, error)
	ClearCart(ctx context.Context, userID string) error

	GetWishlist(ctx context.Context, userID string) ([]Wishlist, error)
	AddToWishlist(ctx context.Context, userID, productID string) error
	RemoveFromWishlist(ctx context.Context, userID, productID string) error
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)
	GetByUserID(ctx context.Context, userID string, page, limit int) ([]Order, int64, error)
	GetAll(ctx context.Context, status string, page, limit int) ([]Order, int64, error)
	UpdateStatus(ctx context.Context, id, status, trackingNumber string) (*Order, error)
	UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) (*Order, error)
}

type AdminRepository interface {
	GetDashboardStats(ctx context.Context) (map[string]interface{}, error)
	GetSettings(ctx context.Context) (map[string]interface{}, error)
	UpdateSetting(ctx context.Context, key string, val interface{}) error
	GetAllReviews(ctx context.Context, status string, page, limit int) ([]Review, int64, error)
	UpdateReviewStatus(ctx context.Context, id, status string) error
}
