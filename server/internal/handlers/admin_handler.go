package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/jaberpatwary/startech/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetDashboardStats(c echo.Context) error {
	var totalUsers, totalOrders, totalProducts, pendingOrders int64
	h.DB.Model(&models.User{}).Where("role = ?", models.RoleUser).Count(&totalUsers)
	h.DB.Model(&models.Order{}).Count(&totalOrders)
	h.DB.Model(&models.Product{}).Count(&totalProducts)
	h.DB.Model(&models.Order{}).Where("status = ?", models.OrderPending).Count(&pendingOrders)

	var totalRevenue struct{ Sum int64 }
	h.DB.Model(&models.Order{}).Where("payment_status = ?", models.PaymentPaid).
		Select("COALESCE(SUM(total), 0) as sum").Scan(&totalRevenue)

	var recentOrders []models.Order
	h.DB.Preload("User").Order("created_at DESC").Limit(5).Find(&recentOrders)

	var lowStock []models.Product
	h.DB.Where("stock < 5 AND status = ?", models.ProductActive).Limit(10).Find(&lowStock)

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success",
		"stats": echo.Map{
			"total_users":    totalUsers,
			"total_orders":   totalOrders,
			"total_products": totalProducts,
			"pending_orders": pendingOrders,
			"total_revenue":  totalRevenue.Sum,
		},
		"recent_orders": recentOrders,
		"low_stock":     lowStock,
	})
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

	q := h.DB.Model(&models.User{})
	if search := c.QueryParam("search"); search != "" {
		q = q.Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	q.Count(&total)

	var users []models.User
	q.Offset((page-1)*limit).Limit(limit).Order("created_at DESC").Find(&users)

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success", "results": users,
		"page": page, "limit": limit, "total_results": total,
		"total_pages": int64(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *AdminHandler) ToggleBlockUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := h.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	user.IsBlocked = !user.IsBlocked
	h.DB.Save(&user)
	msg := "User unblocked"
	if user.IsBlocked {
		msg = "User blocked"
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": msg, "is_blocked": user.IsBlocked})
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

	q := h.DB.Model(&models.Review{}).Preload("User")
	if status := c.QueryParam("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var reviews []models.Review
	q.Offset((page-1)*limit).Limit(limit).Order("created_at DESC").Find(&reviews)

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success", "results": reviews,
		"total_results": total,
	})
}

func (h *AdminHandler) UpdateReviewStatus(c echo.Context) error {
	id := c.Param("id")
	var review models.Review
	if err := h.DB.Where("id = ?", id).First(&review).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Review not found")
	}
	var input struct {
		Status string `json:"status"`
	}
	c.Bind(&input)
	review.Status = input.Status
	h.DB.Save(&review)

	// Recalculate product rating
	if review.Status == models.ReviewApproved {
		var avg struct {
			Avg   float64
			Count int64
		}
		h.DB.Model(&models.Review{}).
			Where("product_id = ? AND status = ?", review.ProductID, models.ReviewApproved).
			Select("AVG(rating) as avg, COUNT(*) as count").Scan(&avg)
		h.DB.Model(&models.Product{}).Where("id = ?", review.ProductID).
			Updates(map[string]interface{}{"rating_avg": avg.Avg, "rating_count": avg.Count})
	}

	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "review": review})
}

func (h *AdminHandler) GetInventory(c echo.Context) error {
	var products []models.Product
	h.DB.Preload("Category").Preload("Brand").Order("stock ASC").Find(&products)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "products": products})
}
