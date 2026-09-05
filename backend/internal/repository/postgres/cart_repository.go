package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/domain"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.WithContext(ctx).
		Preload("Items.Product.Images").
		Preload("Items.Product.Brand").
		Preload("Items.Product.Category").
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create empty cart
			cart = domain.Cart{
				ID:        uuid.NewString(),
				UserID:    userID,
				UpdatedAt: time.Now(),
			}
			if createErr := r.db.WithContext(ctx).Create(&cart).Error; createErr != nil {
				return nil, createErr
			}
			return &cart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, userID, productID string, quantity int) (*domain.Cart, error) {
	cart, err := r.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var existingItem domain.CartItem
	err = r.db.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&existingItem).Error
	if err == nil {
		existingItem.Quantity += quantity
		if existingItem.Quantity < 1 {
			existingItem.Quantity = 1
		}
		if err := r.db.WithContext(ctx).Save(&existingItem).Error; err != nil {
			return nil, err
		}
	} else {
		item := domain.CartItem{
			ID:        uuid.NewString(),
			CartID:    cart.ID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
			return nil, err
		}
	}

	r.db.WithContext(ctx).Model(cart).Update("updated_at", time.Now())
	return r.GetCartByUserID(ctx, userID)
}

func (r *cartRepository) UpdateItem(ctx context.Context, userID, itemID string, quantity int) (*domain.Cart, error) {
	cart, err := r.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if quantity <= 0 {
		return r.RemoveItem(ctx, userID, itemID)
	}

	if err := r.db.WithContext(ctx).Model(&domain.CartItem{}).
		Where("id = ? AND cart_id = ?", itemID, cart.ID).
		Update("quantity", quantity).Error; err != nil {
		return nil, err
	}

	return r.GetCartByUserID(ctx, userID)
}

func (r *cartRepository) RemoveItem(ctx context.Context, userID, itemID string) (*domain.Cart, error) {
	cart, err := r.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&domain.CartItem{}).Error; err != nil {
		return nil, err
	}

	return r.GetCartByUserID(ctx, userID)
}

func (r *cartRepository) ClearCart(ctx context.Context, userID string) error {
	var cart domain.Cart
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error
}

// Wishlist
func (r *cartRepository) GetWishlist(ctx context.Context, userID string) ([]domain.Wishlist, error) {
	var items []domain.Wishlist
	err := r.db.WithContext(ctx).
		Preload("Product.Images").
		Preload("Product.Brand").
		Preload("Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *cartRepository) AddToWishlist(ctx context.Context, userID, productID string) error {
	var existing domain.Wishlist
	if err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&existing).Error; err == nil {
		return nil // already in wishlist
	}
	item := domain.Wishlist{
		ID:        uuid.NewString(),
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *cartRepository) RemoveFromWishlist(ctx context.Context, userID, productID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.Wishlist{}).Error
}
