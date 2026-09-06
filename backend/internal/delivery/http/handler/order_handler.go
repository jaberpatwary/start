package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req struct {
		ShippingName     string `json:"shipping_name"`
		ShippingPhone    string `json:"shipping_phone"`
		ShippingDivision string `json:"shipping_division"`
		ShippingDistrict string `json:"shipping_district"`
		ShippingThana    string `json:"shipping_thana"`
		ShippingAddress  string `json:"shipping_address"`
		ShippingPostal   string `json:"shipping_postal"`
		PaymentMethod    string `json:"payment_method"`
		CouponCode       string `json:"coupon_code"`
		Note             string `json:"note"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	input := usecase.CreateOrderInput{
		UserID:           userID,
		ShippingName:     req.ShippingName,
		ShippingPhone:    req.ShippingPhone,
		ShippingDivision: req.ShippingDivision,
		ShippingDistrict: req.ShippingDistrict,
		ShippingThana:    req.ShippingThana,
		ShippingAddress:  req.ShippingAddress,
		ShippingPostal:   req.ShippingPostal,
		PaymentMethod:    req.PaymentMethod,
		CouponCode:       req.CouponCode,
		Note:             req.Note,
	}

	order, err := h.orderUsecase.CreateOrder(c.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"status":  "success",
		"message": "Order placed successfully",
		"order":   order,
	})
}

func (h *OrderHandler) GetMyOrders(c echo.Context) error {
	userID := c.Get("user_id").(string)
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := h.orderUsecase.GetUserOrders(c.Request().Context(), userID, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch orders")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":          200,
		"status":        "success",
		"orders":        orders,
		"total_results": total,
		"total_pages":   int64(math.Ceil(float64(total) / float64(limit))),
		"page":          page,
		"limit":         limit,
	})
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	id := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	role, _ := c.Get("role").(string)

	order, err := h.orderUsecase.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found")
	}

	// IDOR Protection: User can only view their own order unless they are an admin
	if role != domain.RoleAdmin && order.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view this order")
	}

	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "order": order})
}

func (h *OrderHandler) TrackOrderByNumber(c echo.Context) error {
	number := c.Param("number")
	order, err := h.orderUsecase.TrackOrder(c.Request().Context(), number)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "order": order})
}

func (h *OrderHandler) AdminGetAllOrders(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	status := c.QueryParam("status")

	orders, total, err := h.orderUsecase.GetAllOrders(c.Request().Context(), status, page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch orders")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"code":          200,
		"status":        "success",
		"orders":        orders,
		"total_results": total,
		"total_pages":   int64(math.Ceil(float64(total) / float64(limit))),
		"page":          page,
		"limit":         limit,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	order, err := h.orderUsecase.UpdateOrderStatus(c.Request().Context(), id, req.Status, req.TrackingNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Order status updated",
		"order":   order,
	})
}

func (h *OrderHandler) UpdatePaymentStatus(c echo.Context) error {
	id := c.Param("id")
	var req struct {
		PaymentStatus string `json:"payment_status"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	order, err := h.orderUsecase.UpdatePaymentStatus(c.Request().Context(), id, req.PaymentStatus)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"status":  "success",
		"message": "Payment status updated",
		"order":   order,
	})
}
