package service

import (
	"app/src/model"
	"app/src/response"
	"app/src/validation"
	"encoding/json"
	"errors"
	"fmt"

	"math"

	"math/rand"

	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(userID string, req *validation.CreateOrderRequest) (*model.Order, error)
	GetUserOrders(userID string) ([]model.Order, error)
	GetOrderByID(userID string, role string, orderID string) (*model.Order, error)
	AdminGetAllOrders(page int, limit int, status string) (*response.SuccessWithPaginate[model.Order], error)
	UpdateOrderStatus(orderID string, req *validation.UpdateOrderStatusRequest) (*model.Order, error)
}

type orderServiceImpl struct {
	db                  *gorm.DB
	cartWishlistService CartWishlistService
}

func NewOrderService(db *gorm.DB, cartWishlistService CartWishlistService) OrderService {
	return &orderServiceImpl{
		db:                  db,
		cartWishlistService: cartWishlistService,
	}
}

func (s *orderServiceImpl) CreateOrder(userID string, req *validation.CreateOrderRequest) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must contain at least one item")
	}

	var orderItems []model.OrderItem
	var subtotal int

	for _, itemReq := range req.Items {
		var product model.Product
		if err := s.db.Preload("Images").Where("id = ?", itemReq.ProductID).First(&product).Error; err != nil {
			return nil, fmt.Errorf("product %s not found", itemReq.ProductID)
		}

		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		itemPrice := product.Price
		if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
			itemPrice = *product.DiscountPrice
		}

		imgUrl := ""
		if len(product.Images) > 0 {
			imgUrl = product.Images[0].URL
		}

		orderItems = append(orderItems, model.OrderItem{
			ID:        uuid.New().String(),
			ProductID: product.ID,
			Name:      product.Name,
			Price:     itemPrice,
			Quantity:  itemReq.Quantity,
			Image:     imgUrl,
		})

		subtotal += itemPrice * itemReq.Quantity

		// Deduct stock
		product.Stock -= itemReq.Quantity
		product.Sold += itemReq.Quantity
		s.db.Save(&product)
	}

	shippingFee := 60 // Default Inside Dhaka shipping fee
	if req.ShippingDistrict != "" && req.ShippingDistrict != "Dhaka" {
		shippingFee = 120 // Outside Dhaka
	}

	discount := 0
	if req.CouponCode != "" {
		var coupon model.Coupon
		if err := s.db.Where("code = ? AND active = ?", req.CouponCode, true).First(&coupon).Error; err == nil {
			if coupon.Type == model.CouponFixed {
				discount = coupon.Value
			} else if coupon.Type == model.CouponPercent {
				discount = (subtotal * coupon.Value) / 100
			}
			coupon.UsedCount++
			s.db.Save(&coupon)
		}
	}

	total := subtotal + shippingFee - discount
	if total < 0 {
		total = 0
	}

	orderNum := fmt.Sprintf("ST-%d-%04d", time.Now().Unix(), rand.Intn(10000))

	historyJson, _ := json.Marshal([]map[string]interface{}{
		{
			"status":    model.OrderPending,
			"timestamp": time.Now().Format(time.RFC3339),
			"note":      "Order placed successfully",
		},
	})

	order := model.Order{
		ID:               uuid.New().String(),
		OrderNumber:      orderNum,
		UserID:           userID,
		Items:            orderItems,
		ShippingName:     req.ShippingName,
		ShippingPhone:    req.ShippingPhone,
		ShippingDivision: req.ShippingDivision,
		ShippingDistrict: req.ShippingDistrict,
		ShippingThana:    req.ShippingThana,
		ShippingAddress:  req.ShippingAddress,
		ShippingPostal:   req.ShippingPostal,
		PaymentMethod:    req.PaymentMethod,
		PaymentStatus:    model.PaymentUnpaid,
		Subtotal:         subtotal,
		Discount:         discount,
		ShippingFee:      shippingFee,
		Total:            total,
		Status:           model.OrderPending,
		CouponCode:       req.CouponCode,
		Note:             req.Note,
		StatusHistory:    datatypes.JSON(historyJson),
	}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}

	// Clear user cart after successful checkout
	_ = s.cartWishlistService.ClearCart(userID)

	return s.GetOrderByID(userID, model.RoleUser, order.ID)
}

func (s *orderServiceImpl) GetUserOrders(userID string) ([]model.Order, error) {
	var orders []model.Order
	err := s.db.Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (s *orderServiceImpl) GetOrderByID(userID string, role string, orderID string) (*model.Order, error) {
	var order model.Order
	query := s.db.Preload("Items").Preload("User")

	if role != model.RoleAdmin {
		query = query.Where("id = ? AND user_id = ?", orderID, userID)
	} else {
		query = query.Where("id = ?", orderID)
	}

	if err := query.First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}

	return &order, nil
}

func (s *orderServiceImpl) AdminGetAllOrders(page int, limit int, status string) (*response.SuccessWithPaginate[model.Order], error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	var orders []model.Order
	var totalResults int64

	dbQuery := s.db.Model(&model.Order{}).Preload("Items").Preload("User")
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}

	dbQuery.Count(&totalResults)

	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(limit)))

	return &response.SuccessWithPaginate[model.Order]{
		Code:         200,
		Status:       "success",
		Message:      "Orders fetched successfully",
		Results:      orders,
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(orderID string, req *validation.UpdateOrderStatusRequest) (*model.Order, error) {
	var order model.Order
	if err := s.db.Where("id = ?", orderID).First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if req.Status != "" {
		order.Status = req.Status
	}
	if req.PaymentStatus != "" {
		order.PaymentStatus = req.PaymentStatus
	}

	if err := s.db.Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
