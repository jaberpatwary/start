package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
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

	filter := domain.ProductFilter{
		Page:     page,
		Limit:    limit,
		Category: c.QueryParam("category"),
		Brand:    c.QueryParam("brand"),
		Search:   c.QueryParam("search"),
		Sort:     c.QueryParam("sort"),
		Status:   c.QueryParam("status"),
	}
	if featuredParam := c.QueryParam("featured"); featuredParam != "" {
		featured := featuredParam == "true"
		filter.Featured = &featured
	}

	if minP := c.QueryParam("min_price"); minP != "" {
		filter.MinPrice, _ = strconv.Atoi(minP)
	}
	if maxP := c.QueryParam("max_price"); maxP != "" {
		filter.MaxPrice, _ = strconv.Atoi(maxP)
	}

	products, total, err := h.productUsecase.GetAllProducts(c.Request().Context(), filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch products")
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(http.StatusOK, echo.Map{
		"code":          200,
		"status":        "success",
		"results":       products,
		"page":          page,
		"limit":         limit,
		"total_results": total,
		"total_pages":   totalPages,
	})
}

func (h *ProductHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	product, err := h.productUsecase.GetProductByID(c.Request().Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch product")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "product": product})
}

func (h *ProductHandler) GetBySlug(c echo.Context) error {
	s := c.Param("slug")
	product, err := h.productUsecase.GetProductBySlug(c.Request().Context(), s)
	if err != nil {
		if err == domain.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch product")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "product": product})
}

type createProductRequest struct {
	Name             string   `json:"name"`
	CategoryID       string   `json:"category_id"`
	BrandID          string   `json:"brand_id"`
	Price            int      `json:"price"`
	DiscountPrice    *int     `json:"discount_price"`
	Stock            int      `json:"stock"`
	SKU              string   `json:"sku"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	Tags             []string `json:"tags"`
	Featured         bool     `json:"featured"`
	Status           string   `json:"status"`
	Images           []string `json:"images"`
}

func (h *ProductHandler) Create(c echo.Context) error {
	var req createProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	product := &domain.Product{
		Name:             req.Name,
		CategoryID:       req.CategoryID,
		BrandID:          req.BrandID,
		Price:            req.Price,
		DiscountPrice:    req.DiscountPrice,
		Stock:            req.Stock,
		SKU:              req.SKU,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		Featured:         req.Featured,
		Status:           req.Status,
	}
	if len(req.Tags) > 0 {
		product.Tags = req.Tags
	}

	if err := h.productUsecase.CreateProduct(c.Request().Context(), product, req.Images); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"status":  "success",
		"message": "Product created successfully",
		"product": product,
	})
}

func (h *ProductHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req createProductRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	product := &domain.Product{
		ID:               id,
		Name:             req.Name,
		CategoryID:       req.CategoryID,
		BrandID:          req.BrandID,
		Price:            req.Price,
		DiscountPrice:    req.DiscountPrice,
		Stock:            req.Stock,
		SKU:              req.SKU,
		ShortDescription: req.ShortDescription,
		Description:      req.Description,
		Featured:         req.Featured,
		Status:           req.Status,
	}
	if len(req.Tags) > 0 {
		product.Tags = req.Tags
	}

	if err := h.productUsecase.UpdateProduct(c.Request().Context(), product, req.Images); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Product updated successfully",
		"product": product,
	})
}

func (h *ProductHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.productUsecase.DeleteProduct(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete product")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Product deleted"})
}

func (h *ProductHandler) GetReviews(c echo.Context) error {
	productID := c.Param("id")
	reviews, err := h.productUsecase.GetReviews(c.Request().Context(), productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch reviews")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "reviews": reviews})
}

func (h *ProductHandler) AddReview(c echo.Context) error {
	userID := c.Get("user_id").(string)
	productID := c.Param("id")

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if req.Rating < 1 || req.Rating > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "Rating must be between 1 and 5 stars")
	}

	cleanComment := strings.TrimSpace(req.Comment)
	if len(cleanComment) > 1000 {
		return echo.NewHTTPError(http.StatusBadRequest, "Comment cannot exceed 1000 characters")
	}

	if err := h.productUsecase.CreateReview(c.Request().Context(), userID, productID, req.Rating, cleanComment); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to submit review")
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"status":  "success",
		"message": "Review submitted and pending approval",
	})
}
