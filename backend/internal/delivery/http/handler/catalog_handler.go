package handler

import (
	"net/http"
	"strconv"

	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	catalogUsecase usecase.CatalogUsecase
}

func NewCatalogHandler(catalogUsecase usecase.CatalogUsecase) *CatalogHandler {
	return &CatalogHandler{catalogUsecase: catalogUsecase}
}

// --- Categories ---

func (h *CatalogHandler) GetCategories(c echo.Context) error {
	cats, err := h.catalogUsecase.GetCategories(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch categories")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "categories": cats})
}

func (h *CatalogHandler) CreateCategory(c echo.Context) error {
	var req struct {
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	cat, err := h.catalogUsecase.CreateCategory(c.Request().Context(), req.Name, req.Icon, req.SortOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "category": cat})
}

func (h *CatalogHandler) UpdateCategory(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	cat, err := h.catalogUsecase.UpdateCategory(c.Request().Context(), id, req.Name, req.Icon, req.SortOrder)
	if err != nil {
		if err == domain.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "category": cat})
}

func (h *CatalogHandler) DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if err := h.catalogUsecase.DeleteCategory(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete category")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Category deleted"})
}

// --- Brands ---

func (h *CatalogHandler) GetBrands(c echo.Context) error {
	brands, err := h.catalogUsecase.GetBrands(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch brands")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "brands": brands})
}

func (h *CatalogHandler) CreateBrand(c echo.Context) error {
	var req struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	brand, err := h.catalogUsecase.CreateBrand(c.Request().Context(), req.Name, req.Logo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "brand": brand})
}

func (h *CatalogHandler) UpdateBrand(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Name string `json:"name"`
		Logo string `json:"logo"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	brand, err := h.catalogUsecase.UpdateBrand(c.Request().Context(), id, req.Name, req.Logo)
	if err != nil {
		if err == domain.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Brand not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "brand": brand})
}

func (h *CatalogHandler) DeleteBrand(c echo.Context) error {
	id := c.Param("id")
	if err := h.catalogUsecase.DeleteBrand(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete brand")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Brand deleted"})
}

// --- Banners ---

func (h *CatalogHandler) GetBanners(c echo.Context) error {
	banners, err := h.catalogUsecase.GetBanners(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch banners")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "banners": banners})
}

func (h *CatalogHandler) CreateBanner(c echo.Context) error {
	var banner domain.Banner
	if err := c.Bind(&banner); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.catalogUsecase.CreateBanner(c.Request().Context(), &banner); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "banner": banner})
}

func (h *CatalogHandler) DeleteBanner(c echo.Context) error {
	id := c.Param("id")
	if err := h.catalogUsecase.DeleteBanner(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete banner")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Banner deleted"})
}

// --- Coupons ---

func (h *CatalogHandler) ValidateCoupon(c echo.Context) error {
	code := c.QueryParam("code")
	subtotal, _ := strconv.Atoi(c.QueryParam("subtotal"))

	coupon, discount, err := h.catalogUsecase.ValidateCoupon(c.Request().Context(), code, subtotal)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid or expired coupon")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":     200,
		"status":   "success",
		"coupon":   coupon,
		"discount": discount,
	})
}

func (h *CatalogHandler) CreateCoupon(c echo.Context) error {
	var coupon domain.Coupon
	if err := c.Bind(&coupon); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.catalogUsecase.CreateCoupon(c.Request().Context(), &coupon); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "coupon": coupon})
}

func (h *CatalogHandler) DeleteCoupon(c echo.Context) error {
	id := c.Param("id")
	if err := h.catalogUsecase.DeleteCoupon(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete coupon")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Coupon deleted"})
}
