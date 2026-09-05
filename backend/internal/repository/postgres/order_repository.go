package postgres

import (
	"context"

	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Verify stock & deduct
		for _, item := range order.Items {
			var product domain.Product
			if err := tx.First(&product, "id = ?", item.ProductID).Error; err != nil {
				return err
			}
			if product.Stock < item.Quantity {
				return domain.ErrOutOfStock
			}
			if err := tx.Model(&product).UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
			if err := tx.Model(&product).UpdateColumn("sold", gorm.Expr("sold + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Create order and line items
		return tx.Create(order).Error
	})
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items.Product.Images").
		Preload("User").
		Where("id = ?", id).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).
		Preload("Items.Product.Images").
		Where("order_number = ? OR tracking_number = ?", orderNumber, orderNumber).
		First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Order{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("Items.Product.Images").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) GetAll(ctx context.Context, status string, page, limit int) ([]domain.Order, int64, error) {
	var orders []domain.Order
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("Items").Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id, status, trackingNumber string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	order.Status = status
	if trackingNumber != "" {
		order.TrackingNumber = trackingNumber
	}
	if err := r.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	order.PaymentStatus = paymentStatus
	if err := r.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
