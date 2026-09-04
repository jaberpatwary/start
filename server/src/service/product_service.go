package service

import (
	"app/src/model"
	"app/src/response"
	"app/src/validation"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(req *validation.CreateProductRequest) (*model.Product, error)
	UpdateProduct(productID string, req *validation.UpdateProductRequest) (*model.Product, error)
	DeleteProduct(productID string) error
	GetBySlug(productSlug string) (*model.Product, error)
	GetByID(productID string) (*model.Product, error)
	GetAll(query *validation.QueryProduct) (*response.SuccessWithPaginate[model.Product], error)
}

type productServiceImpl struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return &productServiceImpl{db: db}
}

func (s *productServiceImpl) CreateProduct(req *validation.CreateProductRequest) (*model.Product, error) {
	productSlug := slug.Make(req.Name)
	var count int64
	s.db.Model(&model.Product{}).Where("slug = ? OR sku = ?", productSlug, req.SKU).Count(&count)
	if count > 0 {
		return nil, errors.New("product with this name or SKU already exists")
	}

	status := req.Status
	if status == "" {
		status = model.ProductActive
	}

	product := model.Product{
		ID:               uuid.New().String(),
		Name:             req.Name,
		Slug:             productSlug,
		SKU:              req.SKU,
		CategoryID:       req.CategoryID,
		BrandID:          req.BrandID,
		Price:            req.Price,
		DiscountPrice:    req.DiscountPrice,
		Stock:            req.Stock,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		Status:           status,
		Featured:         req.Featured,
	}

	if err := s.db.Create(&product).Error; err != nil {
		return nil, err
	}

	// Add images
	for i, imgUrl := range req.Images {
		img := model.ProductImage{
			ID:        uuid.New().String(),
			ProductID: product.ID,
			URL:       imgUrl,
			SortOrder: i,
		}
		s.db.Create(&img)
	}

	return s.GetByID(product.ID)
}

func (s *productServiceImpl) UpdateProduct(productID string, req *validation.UpdateProductRequest) (*model.Product, error) {
	var product model.Product
	if err := s.db.Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if req.Name != "" {
		product.Name = req.Name
		product.Slug = slug.Make(req.Name)
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.DiscountPrice != nil {
		product.DiscountPrice = req.DiscountPrice
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.ShortDescription != "" {
		product.ShortDescription = req.ShortDescription
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	product.Featured = req.Featured

	if err := s.db.Save(&product).Error; err != nil {
		return nil, err
	}

	if len(req.Images) > 0 {
		s.db.Where("product_id = ?", product.ID).Delete(&model.ProductImage{})
		for i, imgUrl := range req.Images {
			img := model.ProductImage{
				ID:        uuid.New().String(),
				ProductID: product.ID,
				URL:       imgUrl,
				SortOrder: i,
			}
			s.db.Create(&img)
		}
	}

	return s.GetByID(product.ID)
}

func (s *productServiceImpl) DeleteProduct(productID string) error {
	result := s.db.Where("id = ?", productID).Delete(&model.Product{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (s *productServiceImpl) GetBySlug(productSlug string) (*model.Product, error) {
	var product model.Product
	if err := s.db.Preload("Category").Preload("Brand").Preload("Images").Preload("Reviews.User").Where("slug = ?", productSlug).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func (s *productServiceImpl) GetByID(productID string) (*model.Product, error) {
	var product model.Product
	if err := s.db.Preload("Category").Preload("Brand").Preload("Images").Preload("Reviews.User").Where("id = ?", productID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func (s *productServiceImpl) GetAll(query *validation.QueryProduct) (*response.SuccessWithPaginate[model.Product], error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 12
	}

	var products []model.Product
	var totalResults int64

	dbQuery := s.db.Model(&model.Product{}).Preload("Category").Preload("Brand").Preload("Images")

	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		dbQuery = dbQuery.Where("name ILIKE ? OR short_description ILIKE ? OR sku ILIKE ?", searchPattern, searchPattern, searchPattern)
	}
	if query.CategoryID != "" {
		dbQuery = dbQuery.Where("category_id = ?", query.CategoryID)
	}
	if query.BrandID != "" {
		dbQuery = dbQuery.Where("brand_id = ?", query.BrandID)
	}
	if query.MinPrice > 0 {
		dbQuery = dbQuery.Where("price >= ?", query.MinPrice)
	}
	if query.MaxPrice > 0 {
		dbQuery = dbQuery.Where("price <= ?", query.MaxPrice)
	}
	if query.Featured {
		dbQuery = dbQuery.Where("featured = ?", true)
	}

	dbQuery.Count(&totalResults)

	offset := (query.Page - 1) * query.Limit

	sortOrder := "created_at desc"
	switch query.SortBy {
	case "price_asc":
		sortOrder = "price asc"
	case "price_desc":
		sortOrder = "price desc"
	case "popular":
		sortOrder = "sold desc"
	case "rating":
		sortOrder = "rating_avg desc"
	}

	if err := dbQuery.Offset(offset).Limit(query.Limit).Order(sortOrder).Find(&products).Error; err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return &response.SuccessWithPaginate[model.Product]{
		Code:         200,
		Status:       "success",
		Message:      "Products fetched successfully",
		Results:      products,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}
