package postgres

import (
	"context"
	"encoding/json"

	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) domain.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	var totalUsers int64
	var totalOrders int64
	var totalProducts int64
	var totalRevenue int64
	var lowStockCount int64

	r.db.WithContext(ctx).Model(&domain.User{}).Count(&totalUsers)
	r.db.WithContext(ctx).Model(&domain.Order{}).Count(&totalOrders)
	r.db.WithContext(ctx).Model(&domain.Product{}).Count(&totalProducts)
	r.db.WithContext(ctx).Model(&domain.Product{}).Where("stock < 5").Count(&lowStockCount)
	r.db.WithContext(ctx).Model(&domain.Order{}).Where("payment_status = ?", domain.PaymentPaid).Select("COALESCE(SUM(total), 0)").Scan(&totalRevenue)

	var recentOrders []domain.Order
	r.db.WithContext(ctx).Preload("User").Order("created_at DESC").Limit(5).Find(&recentOrders)

	return map[string]interface{}{
		"total_users":     totalUsers,
		"total_orders":    totalOrders,
		"total_products":  totalProducts,
		"total_revenue":   totalRevenue,
		"low_stock_count": lowStockCount,
		"recent_orders":   recentOrders,
	}, nil
}

func (r *adminRepository) GetSettings(ctx context.Context) (map[string]interface{}, error) {
	var settings []domain.Setting
	if err := r.db.WithContext(ctx).Find(&settings).Error; err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	for _, s := range settings {
		var val interface{}
		_ = json.Unmarshal(s.Value, &val)
		result[s.Key] = val
	}
	return result, nil
}

func (r *adminRepository) UpdateSetting(ctx context.Context, key string, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	setting := domain.Setting{Key: key, Value: datatypes.JSON(bytes)}
	return r.db.WithContext(ctx).Save(&setting).Error
}

func (r *adminRepository) GetAllReviews(ctx context.Context, status string, page, limit int) ([]domain.Review, int64, error) {
	var reviews []domain.Review
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Review{})
	if status != "" {
		q = q.Where("status = ?", status)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *adminRepository) UpdateReviewStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&domain.Review{}).Where("id = ?", id).Update("status", status).Error
}
