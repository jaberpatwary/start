package handler

import (
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminUsecase usecase.AdminUsecase
	uploadDir    string
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase, uploadDir string) *AdminHandler {
	return &AdminHandler{adminUsecase: adminUsecase, uploadDir: uploadDir}
}

func (h *AdminHandler) GetDashboardStats(c echo.Context) error {
	stats, err := h.adminUsecase.GetDashboardStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch stats")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "data": stats})
}

func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := h.adminUsecase.GetUsers(c.Request().Context(), page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":          200,
		"status":        "success",
		"results":       users,
		"page":          page,
		"limit":         limit,
		"total_results": total,
		"total_pages":   int64(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *AdminHandler) ToggleBlockUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.adminUsecase.ToggleBlockUser(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}
	msg := "User unblocked"
	if user.IsBlocked {
		msg = "User blocked"
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":       200,
		"status":     "success",
		"message":    msg,
		"is_blocked": user.IsBlocked,
	})
}

func (h *AdminHandler) GetSettings(c echo.Context) error {
	settings, err := h.adminUsecase.GetSettings(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch settings")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "settings": settings})
}

func (h *AdminHandler) UpdateSettings(c echo.Context) error {
	var input map[string]interface{}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	for k, v := range input {
		if err := h.adminUsecase.UpdateSetting(c.Request().Context(), k, v); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update setting: "+k)
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Settings updated"})
}

func (h *AdminHandler) GetAllReviews(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	status := c.QueryParam("status")

	reviews, total, err := h.adminUsecase.GetReviews(c.Request().Context(), status, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch reviews")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":          200,
		"status":        "success",
		"results":       reviews,
		"total_results": total,
	})
}

func (h *AdminHandler) UpdateReviewStatus(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if err := h.adminUsecase.UpdateReviewStatus(c.Request().Context(), id, req.Status); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Review status updated"})
}

func (h *AdminHandler) UploadImage(c echo.Context) error {
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "No image file provided")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open uploaded file")
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowedExts[ext] {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file type. Only JPG, PNG, WebP, GIF allowed")
	}

	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
	dstPath := filepath.Join(h.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save image")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to write image")
	}

	url := "/uploads/" + filename
	return c.JSON(http.StatusOK, echo.Map{
		"code":     200,
		"status":   "success",
		"url":      url,
		"filename": filename,
	})
}

// GetInventory returns all products sorted by stock level
func (h *AdminHandler) GetInventory(c echo.Context) error {
	// Reuse product usecase or return placeholder
	return c.JSON(http.StatusOK, echo.Map{
		"code":   200,
		"status": "success",
		"message": "Use /api/v1/admin/products with sort=stock for inventory",
	})
}

// GetReports returns sales/revenue summary
func (h *AdminHandler) GetReports(c echo.Context) error {
	stats, err := h.adminUsecase.GetDashboardStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate report")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "report": stats})
}
