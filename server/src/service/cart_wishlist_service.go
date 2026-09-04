package service

import (
	"app/src/model"
	"app/src/validation"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartWishlistService interface {
	GetCart(userID string) (*model.Cart, error)
	AddToCart(userID string, req *validation.AddToCartRequest) (*model.Cart, error)
	UpdateCartItem(userID string, itemID string, req *validation.UpdateCartItemRequest) (*model.Cart, error)
	RemoveCartItem(userID string, itemID string) (*model.Cart, error)
	ClearCart(userID string) error

	GetWishlist(userID string) ([]model.Wishlist, error)
	ToggleWishlist(userID string, productID string) (bool, error)
}

type cartWishlistServiceImpl struct {
	db *gorm.DB
}

func NewCartWishlistService(db *gorm.DB) CartWishlistService {
	return &cartWishlistServiceImpl{db: db}
}

func (s *cartWishlistServiceImpl) GetCart(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := s.db.Preload("Items.Product.Images").Preload("Items.Product.Category").Preload("Items.Product.Brand").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = model.Cart{
				ID:        uuid.New().String(),
				UserID:    userID,
				Items:     []model.CartItem{},
				UpdatedAt: time.Now(),
			}
			s.db.Create(&cart)
			return &cart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (s *cartWishlistServiceImpl) AddToCart(userID string, req *validation.AddToCartRequest) (*model.Cart, error) {
	cart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := s.db.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	var existingItem model.CartItem
	err = s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&existingItem).Error
	if err == nil {
		existingItem.Quantity += req.Quantity
		s.db.Save(&existingItem)
	} else {
		newItem := model.CartItem{
			ID:        uuid.New().String(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		s.db.Create(&newItem)
	}

	return s.GetCart(userID)
}

func (s *cartWishlistServiceImpl) UpdateCartItem(userID string, itemID string, req *validation.UpdateCartItemRequest) (*model.Cart, error) {
	cart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}

	var item model.CartItem
	if err := s.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&item).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	item.Quantity = req.Quantity
	s.db.Save(&item)

	return s.GetCart(userID)
}

func (s *cartWishlistServiceImpl) RemoveCartItem(userID string, itemID string) (*model.Cart, error) {
	cart, err := s.GetCart(userID)
	if err != nil {
		return nil, err
	}

	s.db.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&model.CartItem{})
	return s.GetCart(userID)
}

func (s *cartWishlistServiceImpl) ClearCart(userID string) error {
	var cart model.Cart
	if err := s.db.Where("user_id = ?", userID).First(&cart).Error; err == nil {
		s.db.Where("cart_id = ?", cart.ID).Delete(&model.CartItem{})
	}
	return nil
}

func (s *cartWishlistServiceImpl) GetWishlist(userID string) ([]model.Wishlist, error) {
	var items []model.Wishlist
	err := s.db.Preload("Product.Images").Preload("Product.Category").Preload("Product.Brand").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (s *cartWishlistServiceImpl) ToggleWishlist(userID string, productID string) (bool, error) {
	var existing model.Wishlist
	err := s.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&existing).Error

	if err == nil {
		// Item exists -> remove it
		s.db.Delete(&existing)
		return false, nil
	}

	// Item does not exist -> add it
	newItem := model.Wishlist{
		ID:        uuid.New().String(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	}
	if err := s.db.Create(&newItem).Error; err != nil {
		return false, err
	}

	return true, nil
}
