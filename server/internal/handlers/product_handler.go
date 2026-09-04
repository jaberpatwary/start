package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/jaberpatwary/startech/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (h *ProductHandler) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}

	q := h.DB.Model(&models.Product{}).Preload("Category").Preload("Brand").Preload("Images")

	if search := c.QueryParam("search"); search != "" {
		q = q.Where("name ILIKE ? OR short_description ILIKE ? OR sku ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if catID := c.QueryParam("category_id"); catID != "" {
		q = q.Where("category_id = ?", catID)
	}
	if catSlug := c.QueryParam("category"); catSlug != "" {
		var cat models.Category
		h.DB.Where("slug = ?", catSlug).First(&cat)
		if cat.ID != "" {
			q = q.Where("category_id = ?", cat.ID)
		}
	}
	if brandID := c.QueryParam("brand_id"); brandID != "" {
		q = q.Where("brand_id = ?", brandID)
	}
	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		q = q.Where("price >= ?", minPrice)
	}
	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		q = q.Where("price <= ?", maxPrice)
	}
	if featured := c.QueryParam("featured"); featured == "true" {
		q = q.Where("featured = ?", true)
	}
	if inStock := c.QueryParam("in_stock"); inStock == "true" {
		q = q.Where("stock > 0")
	}

	sortBy := c.QueryParam("sort_by")
	switch sortBy {
	case "price_asc":
		q = q.Order("price ASC")
	case "price_desc":
		q = q.Order("price DESC")
	case "popular":
		q = q.Order("sold DESC")
	case "rating":
		q = q.Order("rating_avg DESC")
	default:
		q = q.Order("created_at DESC")
	}

	var total int64
	q.Count(&total)

	var products []models.Product
	offset := (page - 1) * limit
	q.Offset(offset).Limit(limit).Find(&products)

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success",
		"results":       products,
		"page":          page,
		"limit":         limit,
		"total_results": total,
		"total_pages":   int64(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *ProductHandler) GetBySlug(c echo.Context) error {
	var product models.Product
	if err := h.DB.Preload("Category").Preload("Brand").Preload("Images").
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", models.ReviewApproved).Order("created_at DESC").Limit(20)
		}).Preload("Reviews.User").
		Where("slug = ?", c.Param("slug")).First(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "product": product})
}

func (h *ProductHandler) GetByID(c echo.Context) error {
	var product models.Product
	if err := h.DB.Preload("Category").Preload("Brand").Preload("Images").
		Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "product": product})
}

func (h *ProductHandler) Create(c echo.Context) error {
	var input struct {
		Name             string   `json:"name"`
		SKU              string   `json:"sku"`
		CategoryID       string   `json:"category_id"`
		BrandID          string   `json:"brand_id"`
		Price            int      `json:"price"`
		DiscountPrice    *int     `json:"discount_price"`
		Stock            int      `json:"stock"`
		ShortDescription string   `json:"short_description"`
		Description      string   `json:"description"`
		Images           []string `json:"images"`
		Featured         bool     `json:"featured"`
		Status           string   `json:"status"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if input.Status == "" {
		input.Status = models.ProductActive
	}
	product := models.Product{
		ID:               uuid.NewString(),
		Name:             input.Name,
		Slug:             slug.Make(input.Name),
		SKU:              input.SKU,
		CategoryID:       input.CategoryID,
		BrandID:          input.BrandID,
		Price:            input.Price,
		DiscountPrice:    input.DiscountPrice,
		Stock:            input.Stock,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		Status:           input.Status,
		Featured:         input.Featured,
	}
	if err := h.DB.Create(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	for i, url := range input.Images {
		h.DB.Create(&models.ProductImage{ID: uuid.NewString(), ProductID: product.ID, URL: url, SortOrder: i})
	}
	h.DB.Preload("Category").Preload("Brand").Preload("Images").Where("id = ?", product.ID).First(&product)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "product": product})
}

func (h *ProductHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var product models.Product
	if err := h.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	var input struct {
		Name             string   `json:"name"`
		Price            int      `json:"price"`
		DiscountPrice    *int     `json:"discount_price"`
		Stock            int      `json:"stock"`
		ShortDescription string   `json:"short_description"`
		Description      string   `json:"description"`
		Images           []string `json:"images"`
		Featured         bool     `json:"featured"`
		Status           string   `json:"status"`
		CategoryID       string   `json:"category_id"`
		BrandID          string   `json:"brand_id"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if input.Name != "" {
		product.Name = input.Name
		product.Slug = slug.Make(input.Name)
	}
	if input.Price > 0 {
		product.Price = input.Price
	}
	product.DiscountPrice = input.DiscountPrice
	if input.Stock >= 0 {
		product.Stock = input.Stock
	}
	if input.ShortDescription != "" {
		product.ShortDescription = input.ShortDescription
	}
	if input.Description != "" {
		product.Description = input.Description
	}
	if input.Status != "" {
		product.Status = input.Status
	}
	product.Featured = input.Featured
	if input.CategoryID != "" {
		product.CategoryID = input.CategoryID
	}
	if input.BrandID != "" {
		product.BrandID = input.BrandID
	}
	h.DB.Save(&product)
	if len(input.Images) > 0 {
		h.DB.Where("product_id = ?", product.ID).Delete(&models.ProductImage{})
		for i, url := range input.Images {
			h.DB.Create(&models.ProductImage{ID: uuid.NewString(), ProductID: product.ID, URL: url, SortOrder: i})
		}
	}
	h.DB.Preload("Category").Preload("Brand").Preload("Images").Where("id = ?", product.ID).First(&product)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "product": product})
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.DB.Where("id = ?", id).Delete(&models.Product{}).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Product deleted"})
}

func (h *ProductHandler) AddReview(c echo.Context) error {
	userID := c.Get("user_id").(string)
	productID := c.Param("id")
	var input struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if input.Rating < 1 || input.Rating > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Rating must be between 1 and 5")
	}

	var existing models.Review
	if err := h.DB.Where("product_id = ? AND user_id = ?", productID, userID).First(&existing).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "You have already reviewed this product")
	}

	review := models.Review{
		ID:        uuid.NewString(),
		ProductID: productID,
		UserID:    userID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		Status:    models.ReviewPending,
	}
	h.DB.Create(&review)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "message": "Review submitted, pending approval", "review": review})
}

// Category & Brand handlers
type CatalogHandler struct {
	DB *gorm.DB
}

func NewCatalogHandler(db *gorm.DB) *CatalogHandler {
	return &CatalogHandler{DB: db}
}

func (h *CatalogHandler) GetCategories(c echo.Context) error {
	var categories []models.Category
	h.DB.Order("sort_order ASC").Find(&categories)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "categories": categories})
}

func (h *CatalogHandler) CreateCategory(c echo.Context) error {
	var input struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	cat := models.Category{ID: uuid.NewString(), Name: input.Name, Slug: slug.Make(input.Name), Icon: input.Icon}
	h.DB.Create(&cat)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "category": cat})
}

func (h *CatalogHandler) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var cat models.Category
	if err := h.DB.Where("id = ?", id).First(&cat).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found")
	}
	c.Bind(&cat)
	h.DB.Save(&cat)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "category": cat})
}

func (h *CatalogHandler) DeleteCategory(c echo.Context) error {
	h.DB.Where("id = ?", c.Param("id")).Delete(&models.Category{})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Category deleted"})
}

func (h *CatalogHandler) GetBrands(c echo.Context) error {
	var brands []models.Brand
	h.DB.Order("name ASC").Find(&brands)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "brands": brands})
}

func (h *CatalogHandler) CreateBrand(c echo.Context) error {
	var input struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	brand := models.Brand{ID: uuid.NewString(), Name: input.Name, Slug: slug.Make(input.Name), Logo: input.Logo}
	h.DB.Create(&brand)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "brand": brand})
}

func (h *CatalogHandler) UpdateBrand(c echo.Context) error {
	id := c.Param("id")
	var brand models.Brand
	if err := h.DB.Where("id = ?", id).First(&brand).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Brand not found")
	}
	c.Bind(&brand)
	h.DB.Save(&brand)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "brand": brand})
}

func (h *CatalogHandler) DeleteBrand(c echo.Context) error {
	h.DB.Where("id = ?", c.Param("id")).Delete(&models.Brand{})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Brand deleted"})
}

func (h *CatalogHandler) GetBanners(c echo.Context) error {
	var banners []models.Banner
	h.DB.Where("active = ?", true).Order("sort_order ASC").Find(&banners)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "banners": banners})
}

func (h *CatalogHandler) AdminGetBanners(c echo.Context) error {
	var banners []models.Banner
	h.DB.Order("sort_order ASC").Find(&banners)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "banners": banners})
}

func (h *CatalogHandler) CreateBanner(c echo.Context) error {
	var banner models.Banner
	c.Bind(&banner)
	banner.ID = uuid.NewString()
	h.DB.Create(&banner)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "banner": banner})
}

func (h *CatalogHandler) UpdateBanner(c echo.Context) error {
	id := c.Param("id")
	var banner models.Banner
	if err := h.DB.Where("id = ?", id).First(&banner).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Banner not found")
	}
	c.Bind(&banner)
	h.DB.Save(&banner)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "banner": banner})
}

func (h *CatalogHandler) DeleteBanner(c echo.Context) error {
	h.DB.Where("id = ?", c.Param("id")).Delete(&models.Banner{})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Banner deleted"})
}

func (h *CatalogHandler) GetCoupons(c echo.Context) error {
	var coupons []models.Coupon
	h.DB.Find(&coupons)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "coupons": coupons})
}

func (h *CatalogHandler) CreateCoupon(c echo.Context) error {
	var coupon models.Coupon
	c.Bind(&coupon)
	coupon.ID = uuid.NewString()
	h.DB.Create(&coupon)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "coupon": coupon})
}

func (h *CatalogHandler) UpdateCoupon(c echo.Context) error {
	id := c.Param("id")
	var coupon models.Coupon
	if err := h.DB.Where("id = ?", id).First(&coupon).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Coupon not found")
	}
	c.Bind(&coupon)
	h.DB.Save(&coupon)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "coupon": coupon})
}

func (h *CatalogHandler) DeleteCoupon(c echo.Context) error {
	h.DB.Where("id = ?", c.Param("id")).Delete(&models.Coupon{})
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Coupon deleted"})
}

func (h *CatalogHandler) ValidateCoupon(c echo.Context) error {
	code := c.QueryParam("code")
	var coupon models.Coupon
	if err := h.DB.Where("code = ? AND active = ?", code, true).First(&coupon).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Coupon not found or expired")
	}
	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return echo.NewHTTPError(http.StatusGone, "Coupon usage limit reached")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "coupon": coupon})
}
