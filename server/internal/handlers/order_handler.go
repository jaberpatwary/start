package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB          *gorm.DB
	CartHandler *CartHandler
}

func NewOrderHandler(db *gorm.DB, cartHandler *CartHandler) *OrderHandler {
	return &OrderHandler{DB: db, CartHandler: cartHandler}
}

type CreateOrderInput struct {
	Items            []OrderItemInput `json:"items"`
	ShippingName     string           `json:"shipping_name"`
	ShippingPhone    string           `json:"shipping_phone"`
	ShippingDivision string           `json:"shipping_division"`
	ShippingDistrict string           `json:"shipping_district"`
	ShippingThana    string           `json:"shipping_thana"`
	ShippingAddress  string           `json:"shipping_address"`
	ShippingPostal   string           `json:"shipping_postal"`
	PaymentMethod    string           `json:"payment_method"`
	CouponCode       string           `json:"coupon_code"`
	Note             string           `json:"note"`
}

type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var input CreateOrderInput
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if len(input.Items) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Order must have at least one item")
	}

	var orderItems []models.OrderItem
	subtotal := 0

	for _, item := range input.Items {
		var product models.Product
		if err := h.DB.Preload("Images").Where("id = ?", item.ProductID).First(&product).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Product %s not found", item.ProductID))
		}
		if product.Stock < item.Quantity {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Insufficient stock for %s", product.Name))
		}

		price := product.Price
		if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
			price = *product.DiscountPrice
		}

		imgURL := ""
		if len(product.Images) > 0 {
			imgURL = product.Images[0].URL
		}

		orderItems = append(orderItems, models.OrderItem{
			ID: uuid.NewString(), ProductID: product.ID,
			Name: product.Name, Price: price, Quantity: item.Quantity, Image: imgURL,
		})
		subtotal += price * item.Quantity

		product.Stock -= item.Quantity
		product.Sold += item.Quantity
		h.DB.Save(&product)
	}

	shippingFee := 60
	if input.ShippingDistrict != "" && input.ShippingDistrict != "Dhaka" {
		shippingFee = 120
	}

	discount := 0
	couponUsed := ""
	if input.CouponCode != "" {
		var coupon models.Coupon
		if err := h.DB.Where("code = ? AND active = ?", input.CouponCode, true).First(&coupon).Error; err == nil {
			if subtotal >= coupon.MinOrder {
				if coupon.Type == models.CouponPercent {
					discount = (subtotal * coupon.Value) / 100
				} else {
					discount = coupon.Value
				}
				coupon.UsedCount++
				h.DB.Save(&coupon)
				couponUsed = coupon.Code
			}
		}
	}

	total := subtotal + shippingFee - discount
	if total < 0 {
		total = 0
	}

	orderNum := fmt.Sprintf("ST-%d-%04d", time.Now().Unix(), rand.Intn(9999))
	histJSON, _ := json.Marshal([]map[string]interface{}{
		{"status": models.OrderPending, "timestamp": time.Now().Format(time.RFC3339), "note": "Order placed"},
	})

	order := models.Order{
		ID:               uuid.NewString(),
		OrderNumber:      orderNum,
		UserID:           userID,
		Items:            orderItems,
		ShippingName:     input.ShippingName,
		ShippingPhone:    input.ShippingPhone,
		ShippingDivision: input.ShippingDivision,
		ShippingDistrict: input.ShippingDistrict,
		ShippingThana:    input.ShippingThana,
		ShippingAddress:  input.ShippingAddress,
		ShippingPostal:   input.ShippingPostal,
		PaymentMethod:    input.PaymentMethod,
		PaymentStatus:    models.PaymentUnpaid,
		Subtotal:         subtotal,
		Discount:         discount,
		ShippingFee:      shippingFee,
		Total:            total,
		Status:           models.OrderPending,
		CouponCode:       couponUsed,
		Note:             input.Note,
		StatusHistory:    datatypes.JSON(histJSON),
	}

	if err := h.DB.Create(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create order")
	}

	h.CartHandler.ClearCart(userID)

	h.DB.Preload("Items").Where("id = ?", order.ID).First(&order)
	return c.JSON(http.StatusCreated, echo.Map{
		"code": 201, "status": "success", "message": "Order placed successfully", "order": order,
	})
}

func (h *OrderHandler) GetMyOrders(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var orders []models.Order
	h.DB.Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "orders": orders})
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	userID := c.Get("user_id").(string)
	role, _ := c.Get("role").(string)
	orderID := c.Param("id")

	var order models.Order
	query := h.DB.Preload("Items.Product").Preload("User")
	if role != models.RoleAdmin {
		query = query.Where("id = ? AND user_id = ?", orderID, userID)
	} else {
		query = query.Where("id = ?", orderID)
	}
	if err := query.First(&order).Error; err != nil {
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

	q := h.DB.Model(&models.Order{}).Preload("Items").Preload("User")
	if status := c.QueryParam("status"); status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	q.Count(&total)

	var orders []models.Order
	q.Offset((page - 1) * limit).Limit(limit).Order("created_at DESC").Find(&orders)

	return c.JSON(http.StatusOK, echo.Map{
		"code": 200, "status": "success", "results": orders,
		"page": page, "limit": limit, "total_results": total,
		"total_pages": int64(math.Ceil(float64(total) / float64(limit))),
	})
}

func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	orderID := c.Param("id")
	var order models.Order
	if err := h.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found")
	}
	var input struct {
		Status         string `json:"status"`
		PaymentStatus  string `json:"payment_status"`
		TrackingNumber string `json:"tracking_number"`
	}
	c.Bind(&input)
	if input.Status != "" {
		order.Status = input.Status
	}
	if input.PaymentStatus != "" {
		order.PaymentStatus = input.PaymentStatus
	}
	if input.TrackingNumber != "" {
		order.TrackingNumber = input.TrackingNumber
	}
	h.DB.Save(&order)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "order": order})
}

func (h *OrderHandler) TrackOrderByNumber(c echo.Context) error {
	number := c.Param("number")
	var order models.Order
	if err := h.DB.Preload("Items").
		Where("order_number = ? OR tracking_number = ?", number, number).
		First(&order).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Order not found with that order or tracking number")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "order": order})
}

