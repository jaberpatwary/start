package postgres

import (
	"context"

	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/gorm"
)

type catalogRepository struct {
	db *gorm.DB
}

func NewCatalogRepository(db *gorm.DB) domain.CatalogRepository {
	return &catalogRepository{db: db}
}

// Categories
func (r *catalogRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

func (r *catalogRepository) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *catalogRepository) CreateCategory(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *catalogRepository) UpdateCategory(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *catalogRepository) DeleteCategory(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Category{}).Error
}

// Brands
func (r *catalogRepository) GetBrands(ctx context.Context) ([]domain.Brand, error) {
	var brands []domain.Brand
	err := r.db.WithContext(ctx).Order("name ASC").Find(&brands).Error
	return brands, err
}

func (r *catalogRepository) GetBrandByID(ctx context.Context, id string) (*domain.Brand, error) {
	var brand domain.Brand
	err := r.db.WithContext(ctx).First(&brand, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &brand, nil
}

func (r *catalogRepository) CreateBrand(ctx context.Context, brand *domain.Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

func (r *catalogRepository) UpdateBrand(ctx context.Context, brand *domain.Brand) error {
	return r.db.WithContext(ctx).Save(brand).Error
}

func (r *catalogRepository) DeleteBrand(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Brand{}).Error
}

// Banners
func (r *catalogRepository) GetBanners(ctx context.Context) ([]domain.Banner, error) {
	var banners []domain.Banner
	err := r.db.WithContext(ctx).Where("active = true").Order("sort_order ASC").Find(&banners).Error
	return banners, err
}

func (r *catalogRepository) CreateBanner(ctx context.Context, banner *domain.Banner) error {
	return r.db.WithContext(ctx).Create(banner).Error
}

func (r *catalogRepository) DeleteBanner(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Banner{}).Error
}

// Coupons
func (r *catalogRepository) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := r.db.WithContext(ctx).Where("UPPER(code) = UPPER(?) AND active = true", code).First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &coupon, nil
}

func (r *catalogRepository) CreateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

func (r *catalogRepository) DeleteCoupon(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Coupon{}).Error
}

func (r *catalogRepository) IncrementCouponUsage(ctx context.Context, code string) error {
	return r.db.WithContext(ctx).Model(&domain.Coupon{}).
		Where("UPPER(code) = UPPER(?)", code).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}
