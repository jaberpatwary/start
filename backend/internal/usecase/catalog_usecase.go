package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jaberpatwary/startech/internal/domain"
)

type CatalogUsecase interface {
	GetCategories(ctx context.Context) ([]domain.Category, error)
	CreateCategory(ctx context.Context, name, icon string, sortOrder int) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id, name, icon string, sortOrder int) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error

	GetBrands(ctx context.Context) ([]domain.Brand, error)
	CreateBrand(ctx context.Context, name, logo string) (*domain.Brand, error)
	UpdateBrand(ctx context.Context, id, name, logo string) (*domain.Brand, error)
	DeleteBrand(ctx context.Context, id string) error

	GetBanners(ctx context.Context) ([]domain.Banner, error)
	CreateBanner(ctx context.Context, banner *domain.Banner) error
	DeleteBanner(ctx context.Context, id string) error

	ValidateCoupon(ctx context.Context, code string, orderSubtotal int) (*domain.Coupon, int, error)
	CreateCoupon(ctx context.Context, coupon *domain.Coupon) error
	DeleteCoupon(ctx context.Context, id string) error
}

type catalogUsecase struct {
	catalogRepo domain.CatalogRepository
}

func NewCatalogUsecase(catalogRepo domain.CatalogRepository) CatalogUsecase {
	return &catalogUsecase{
		catalogRepo: catalogRepo,
	}
}

func (u *catalogUsecase) GetCategories(ctx context.Context) ([]domain.Category, error) {
	return u.catalogRepo.GetCategories(ctx)
}

func (u *catalogUsecase) CreateCategory(ctx context.Context, name, icon string, sortOrder int) (*domain.Category, error) {
	cat := &domain.Category{
		ID:        uuid.NewString(),
		Name:      name,
		Slug:      slug.Make(name),
		Icon:      icon,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}
	if err := u.catalogRepo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (u *catalogUsecase) UpdateCategory(ctx context.Context, id, name, icon string, sortOrder int) (*domain.Category, error) {
	cat, err := u.catalogRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		cat.Name = name
		cat.Slug = slug.Make(name)
	}
	if icon != "" {
		cat.Icon = icon
	}
	cat.SortOrder = sortOrder
	if err := u.catalogRepo.UpdateCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (u *catalogUsecase) DeleteCategory(ctx context.Context, id string) error {
	return u.catalogRepo.DeleteCategory(ctx, id)
}

func (u *catalogUsecase) GetBrands(ctx context.Context) ([]domain.Brand, error) {
	return u.catalogRepo.GetBrands(ctx)
}

func (u *catalogUsecase) CreateBrand(ctx context.Context, name, logo string) (*domain.Brand, error) {
	brand := &domain.Brand{
		ID:        uuid.NewString(),
		Name:      name,
		Slug:      slug.Make(name),
		Logo:      logo,
		CreatedAt: time.Now(),
	}
	if err := u.catalogRepo.CreateBrand(ctx, brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (u *catalogUsecase) UpdateBrand(ctx context.Context, id, name, logo string) (*domain.Brand, error) {
	brand, err := u.catalogRepo.GetBrandByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != "" {
		brand.Name = name
		brand.Slug = slug.Make(name)
	}
	if logo != "" {
		brand.Logo = logo
	}
	if err := u.catalogRepo.UpdateBrand(ctx, brand); err != nil {
		return nil, err
	}
	return brand, nil
}

func (u *catalogUsecase) DeleteBrand(ctx context.Context, id string) error {
	return u.catalogRepo.DeleteBrand(ctx, id)
}

func (u *catalogUsecase) GetBanners(ctx context.Context) ([]domain.Banner, error) {
	return u.catalogRepo.GetBanners(ctx)
}

func (u *catalogUsecase) CreateBanner(ctx context.Context, banner *domain.Banner) error {
	banner.ID = uuid.NewString()
	banner.CreatedAt = time.Now()
	banner.Active = true
	return u.catalogRepo.CreateBanner(ctx, banner)
}

func (u *catalogUsecase) DeleteBanner(ctx context.Context, id string) error {
	return u.catalogRepo.DeleteBanner(ctx, id)
}

func (u *catalogUsecase) ValidateCoupon(ctx context.Context, code string, orderSubtotal int) (*domain.Coupon, int, error) {
	coupon, err := u.catalogRepo.GetCouponByCode(ctx, code)
	if err != nil {
		return nil, 0, domain.ErrCouponInvalid
	}

	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
		return nil, 0, domain.ErrCouponInvalid
	}

	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return nil, 0, domain.ErrCouponInvalid
	}

	if orderSubtotal < coupon.MinOrder {
		return nil, 0, domain.ErrCouponInvalid
	}

	discount := 0
	if coupon.Type == domain.CouponPercent {
		discount = (orderSubtotal * coupon.Value) / 100
	} else {
		discount = coupon.Value
	}

	if discount > orderSubtotal {
		discount = orderSubtotal
	}

	return coupon, discount, nil
}

func (u *catalogUsecase) CreateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	coupon.ID = uuid.NewString()
	coupon.CreatedAt = time.Now()
	coupon.Active = true
	return u.catalogRepo.CreateCoupon(ctx, coupon)
}

func (u *catalogUsecase) DeleteCoupon(ctx context.Context, id string) error {
	return u.catalogRepo.DeleteCoupon(ctx, id)
}
