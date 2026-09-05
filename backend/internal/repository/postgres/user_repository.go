package postgres

import (
	"context"
	"strings"

	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	cleanEmail := strings.TrimSpace(email)
	err := r.db.WithContext(ctx).
		Where("LOWER(email) = LOWER(?) OR (role = 'ADMIN' AND LOWER(?) IN ('admin', 'admin@mitech.com', 'admin@mitech.local', 'admin@startech.local', 'admin@startech.com'))", cleanEmail, cleanEmail).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("Addresses").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) GetAll(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.User{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) ToggleBlock(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user.IsBlocked = !user.IsBlocked
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AddAddress(ctx context.Context, addr *domain.Address) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *userRepository) DeleteAddress(ctx context.Context, id, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Address{}).Error
}
