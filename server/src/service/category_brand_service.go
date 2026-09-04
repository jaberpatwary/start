package service

import (
	"app/src/model"
	"app/src/validation"
	"errors"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CategoryBrandService interface {
	GetAllCategories() ([]model.Category, error)
	CreateCategory(req *validation.CategoryRequest) (*model.Category, error)
	GetAllBrands() ([]model.Brand, error)
	CreateBrand(req *validation.BrandRequest) (*model.Brand, error)
}

type categoryBrandServiceImpl struct {
	db *gorm.DB
}

func NewCategoryBrandService(db *gorm.DB) CategoryBrandService {
	return &categoryBrandServiceImpl{db: db}
}

func (s *categoryBrandServiceImpl) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := s.db.Order("sort_order asc, name asc").Find(&categories).Error
	return categories, err
}

func (s *categoryBrandServiceImpl) CreateCategory(req *validation.CategoryRequest) (*model.Category, error) {
	categorySlug := slug.Make(req.Name)
	var count int64
	s.db.Model(&model.Category{}).Where("slug = ?", categorySlug).Count(&count)
	if count > 0 {
		return nil, errors.New("category with this name already exists")
	}

	category := model.Category{
		ID:   uuid.New().String(),
		Name: req.Name,
		Slug: categorySlug,
		Icon: req.Icon,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categoryBrandServiceImpl) GetAllBrands() ([]model.Brand, error) {
	var brands []model.Brand
	err := s.db.Order("name asc").Find(&brands).Error
	return brands, err
}

func (s *categoryBrandServiceImpl) CreateBrand(req *validation.BrandRequest) (*model.Brand, error) {
	brandSlug := slug.Make(req.Name)
	var count int64
	s.db.Model(&model.Brand{}).Where("slug = ?", brandSlug).Count(&count)
	if count > 0 {
		return nil, errors.New("brand with this name already exists")
	}

	brand := model.Brand{
		ID:   uuid.New().String(),
		Name: req.Name,
		Slug: brandSlug,
		Logo: req.Logo,
	}

	if err := s.db.Create(&brand).Error; err != nil {
		return nil, err
	}
	return &brand, nil
}
