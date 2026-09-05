package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(ctx context.Context, f domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Product{})

	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	} else {
		q = q.Where("status != ?", domain.ProductDraft)
	}

	if f.Category != "" {
		var cat domain.Category
		if err := r.db.WithContext(ctx).Where("slug = ? OR id = ?", f.Category, f.Category).First(&cat).Error; err == nil {
			q = q.Where("category_id = ?", cat.ID)
		}
	}

	if f.Brand != "" {
		var brand domain.Brand
		if err := r.db.WithContext(ctx).Where("slug = ? OR id = ?", f.Brand, f.Brand).First(&brand).Error; err == nil {
			q = q.Where("brand_id = ?", brand.ID)
		}
	}

	if f.Search != "" {
		searchPattern := "%" + f.Search + "%"
		q = q.Where("name ILIKE ? OR short_description ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if f.MinPrice > 0 {
		q = q.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice > 0 {
		q = q.Where("price <= ?", f.MaxPrice)
	}
	if f.Rating > 0 {
		q = q.Where("rating_avg >= ?", f.Rating)
	}
	if f.Featured != nil {
		q = q.Where("featured = ?", *f.Featured)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch f.Sort {
	case "price_asc":
		q = q.Order("price ASC")
	case "price_desc":
		q = q.Order("price DESC")
	case "rating":
		q = q.Order("rating_avg DESC")
	case "popular":
		q = q.Order("sold DESC")
	default:
		q = q.Order("created_at DESC")
	}

	if f.Limit <= 0 {
		f.Limit = 12
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	err := q.Preload("Images").Preload("Category").Preload("Brand").
		Limit(f.Limit).Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Images").Preload("Category").Preload("Brand").
		Preload("Reviews", "status = ?", domain.ReviewApproved).
		Preload("Reviews.User").
		Where("id = ?", id).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Images").Preload("Category").Preload("Brand").
		Preload("Reviews", "status = ?", domain.ReviewApproved).
		Preload("Reviews.User").
		Where("slug = ?", slug).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product, images []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		for i, url := range images {
			img := domain.ProductImage{
				ID:        uuid.NewString(),
				ProductID: product.ID,
				URL:       url,
				SortOrder: i,
			}
			if err := tx.Create(&img).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product, images []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(product).Error; err != nil {
			return err
		}
		if len(images) > 0 {
			if err := tx.Where("product_id = ?", product.ID).Delete(&domain.ProductImage{}).Error; err != nil {
				return err
			}
			for i, url := range images {
				img := domain.ProductImage{
					ID:        uuid.NewString(),
					ProductID: product.ID,
					URL:       url,
					SortOrder: i,
				}
				if err := tx.Create(&img).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Product{}).Error
}

func (r *productRepository) GetReviews(ctx context.Context, productID string) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("product_id = ? AND status = ?", productID, domain.ReviewApproved).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *productRepository) CreateReview(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}
