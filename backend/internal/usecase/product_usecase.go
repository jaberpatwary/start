package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jaberpatwary/startech/internal/domain"
)

type ProductUsecase interface {
	GetAllProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error)
	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product, images []string) error
	UpdateProduct(ctx context.Context, product *domain.Product, images []string) error
	DeleteProduct(ctx context.Context, id string) error
	GetReviews(ctx context.Context, productID string) ([]domain.Review, error)
	CreateReview(ctx context.Context, userID, productID string, rating int, comment string) error
}

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetAllProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	return u.productRepo.GetAll(ctx, filter)
}

func (u *productUsecase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *productUsecase) GetProductBySlug(ctx context.Context, s string) (*domain.Product, error) {
	return u.productRepo.GetBySlug(ctx, s)
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product, images []string) error {
	product.ID = uuid.NewString()
	if product.Slug == "" {
		product.Slug = fmt.Sprintf("%s-%s", slug.Make(product.Name), product.ID[:8])
	}
	if product.SKU == "" {
		product.SKU = fmt.Sprintf("SKU-%s", product.ID[:8])
	}
	if product.Status == "" {
		product.Status = domain.ProductActive
	}
	return u.productRepo.Create(ctx, product, images)
}

func (u *productUsecase) UpdateProduct(ctx context.Context, product *domain.Product, images []string) error {
	if product.Name != "" && product.Slug == "" {
		product.Slug = slug.Make(product.Name)
	}
	return u.productRepo.Update(ctx, product, images)
}

func (u *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *productUsecase) GetReviews(ctx context.Context, productID string) ([]domain.Review, error) {
	return u.productRepo.GetReviews(ctx, productID)
}

func (u *productUsecase) CreateReview(ctx context.Context, userID, productID string, rating int, comment string) error {
	review := &domain.Review{
		ID:        uuid.NewString(),
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   comment,
		Status:    domain.ReviewPending,
		CreatedAt: time.Now(),
	}
	return u.productRepo.CreateReview(ctx, review)
}
