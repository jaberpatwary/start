package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/datatypes"
)

type CreateOrderInput struct {
	UserID           string
	ShippingName     string
	ShippingPhone    string
	ShippingDivision string
	ShippingDistrict string
	ShippingThana    string
	ShippingAddress  string
	ShippingPostal   string
	PaymentMethod    string
	CouponCode       string
	Note             string
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error)
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	TrackOrder(ctx context.Context, orderOrTrackingNumber string) (*domain.Order, error)
	GetUserOrders(ctx context.Context, userID string, page, limit int) ([]domain.Order, int64, error)
	GetAllOrders(ctx context.Context, status string, page, limit int) ([]domain.Order, int64, error)
	UpdateOrderStatus(ctx context.Context, id, status, trackingNumber string) (*domain.Order, error)
	UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) (*domain.Order, error)
}

type orderUsecase struct {
	orderRepo   domain.OrderRepository
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
	catalogRepo domain.CatalogRepository
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	catalogRepo domain.CatalogRepository,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		catalogRepo: catalogRepo,
	}
}

func (u *orderUsecase) CreateOrder(ctx context.Context, input CreateOrderInput) (*domain.Order, error) {
	cart, err := u.cartRepo.GetCartByUserID(ctx, input.UserID)
	if err != nil || len(cart.Items) == 0 {
		return nil, domain.ErrBadRequest
	}

	// Calculate subtotal and build items
	subtotal := 0
	var orderItems []domain.OrderItem
	orderID := uuid.NewString()

	for _, item := range cart.Items {
		product, err := u.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("%s is out of stock (available: %d)", product.Name, product.Stock)
		}

		price := product.Price
		if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
			price = *product.DiscountPrice
		}
		itemTotal := price * item.Quantity
		subtotal += itemTotal

		img := ""
		if len(product.Images) > 0 {
			img = product.Images[0].URL
		}

		orderItems = append(orderItems, domain.OrderItem{
			ID:        uuid.NewString(),
			OrderID:   orderID,
			ProductID: product.ID,
			Name:      product.Name,
			Price:     price,
			Quantity:  item.Quantity,
			Image:     img,
		})
	}

	// Apply coupon if provided
	discount := 0
	if input.CouponCode != "" {
		coupon, err := u.catalogRepo.GetCouponByCode(ctx, input.CouponCode)
		if err == nil && coupon != nil {
			if coupon.ExpiresAt == nil || !coupon.ExpiresAt.Before(time.Now()) {
				if coupon.MaxUses == 0 || coupon.UsedCount < coupon.MaxUses {
					if subtotal >= coupon.MinOrder {
						if coupon.Type == domain.CouponPercent {
							discount = (subtotal * coupon.Value) / 100
						} else {
							discount = coupon.Value
						}
						_ = u.catalogRepo.IncrementCouponUsage(ctx, input.CouponCode)
					}
				}
			}
		}
	}

	// Calculate shipping fee (e.g. Inside Dhaka 60, Outside Dhaka 120)
	shippingFee := 60
	if input.ShippingDivision != "" && input.ShippingDivision != "Dhaka" {
		shippingFee = 120
	}

	total := subtotal - discount + shippingFee
	if total < 0 {
		total = 0
	}

	orderNumber := fmt.Sprintf("ST-%d%04d", time.Now().Unix()%1000000, rand.Intn(10000))
	trackingNumber := fmt.Sprintf("TRK-%d", time.Now().UnixNano()%1000000000)

	history, _ := json.Marshal([]map[string]interface{}{
		{"status": domain.OrderPending, "time": time.Now(), "note": "Order placed successfully"},
	})

	order := &domain.Order{
		ID:               orderID,
		OrderNumber:      orderNumber,
		UserID:           input.UserID,
		Items:            orderItems,
		ShippingName:     input.ShippingName,
		ShippingPhone:    input.ShippingPhone,
		ShippingDivision: input.ShippingDivision,
		ShippingDistrict: input.ShippingDistrict,
		ShippingThana:    input.ShippingThana,
		ShippingAddress:  input.ShippingAddress,
		ShippingPostal:   input.ShippingPostal,
		PaymentMethod:    input.PaymentMethod,
		PaymentStatus:    domain.PaymentUnpaid,
		Subtotal:         subtotal,
		Discount:         discount,
		ShippingFee:      shippingFee,
		Total:            total,
		Status:           domain.OrderPending,
		TrackingNumber:   trackingNumber,
		CouponCode:       input.CouponCode,
		Note:             input.Note,
		StatusHistory:    datatypes.JSON(history),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := u.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Clear user's cart
	_ = u.cartRepo.ClearCart(ctx, input.UserID)

	return order, nil
}

func (u *orderUsecase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return u.orderRepo.GetByID(ctx, id)
}

func (u *orderUsecase) TrackOrder(ctx context.Context, orderOrTrackingNumber string) (*domain.Order, error) {
	return u.orderRepo.GetByOrderNumber(ctx, orderOrTrackingNumber)
}

func (u *orderUsecase) GetUserOrders(ctx context.Context, userID string, page, limit int) ([]domain.Order, int64, error) {
	return u.orderRepo.GetByUserID(ctx, userID, page, limit)
}

func (u *orderUsecase) GetAllOrders(ctx context.Context, status string, page, limit int) ([]domain.Order, int64, error) {
	return u.orderRepo.GetAll(ctx, status, page, limit)
}

func (u *orderUsecase) UpdateOrderStatus(ctx context.Context, id, status, trackingNumber string) (*domain.Order, error) {
	return u.orderRepo.UpdateStatus(ctx, id, status, trackingNumber)
}

func (u *orderUsecase) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) (*domain.Order, error) {
	return u.orderRepo.UpdatePaymentStatus(ctx, id, paymentStatus)
}
