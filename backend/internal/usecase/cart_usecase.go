package usecase

import (
	"context"

	"github.com/jaberpatwary/startech/internal/domain"
)

type CartUsecase interface {
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddToCart(ctx context.Context, userID, productID string, quantity int) (*domain.Cart, error)
	UpdateItem(ctx context.Context, userID, itemID string, quantity int) (*domain.Cart, error)
	RemoveItem(ctx context.Context, userID, itemID string) (*domain.Cart, error)
	ClearCart(ctx context.Context, userID string) error

	GetWishlist(ctx context.Context, userID string) ([]domain.Wishlist, error)
	AddToWishlist(ctx context.Context, userID, productID string) error
	RemoveFromWishlist(ctx context.Context, userID, productID string) error
}

type cartUsecase struct {
	cartRepo domain.CartRepository
}

func NewCartUsecase(cartRepo domain.CartRepository) CartUsecase {
	return &cartUsecase{
		cartRepo: cartRepo,
	}
}

func (u *cartUsecase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return u.cartRepo.GetCartByUserID(ctx, userID)
}

func (u *cartUsecase) AddToCart(ctx context.Context, userID, productID string, quantity int) (*domain.Cart, error) {
	if quantity <= 0 {
		quantity = 1
	}
	return u.cartRepo.AddItem(ctx, userID, productID, quantity)
}

func (u *cartUsecase) UpdateItem(ctx context.Context, userID, itemID string, quantity int) (*domain.Cart, error) {
	return u.cartRepo.UpdateItem(ctx, userID, itemID, quantity)
}

func (u *cartUsecase) RemoveItem(ctx context.Context, userID, itemID string) (*domain.Cart, error) {
	return u.cartRepo.RemoveItem(ctx, userID, itemID)
}

func (u *cartUsecase) ClearCart(ctx context.Context, userID string) error {
	return u.cartRepo.ClearCart(ctx, userID)
}

func (u *cartUsecase) GetWishlist(ctx context.Context, userID string) ([]domain.Wishlist, error) {
	return u.cartRepo.GetWishlist(ctx, userID)
}

func (u *cartUsecase) AddToWishlist(ctx context.Context, userID, productID string) error {
	return u.cartRepo.AddToWishlist(ctx, userID, productID)
}

func (u *cartUsecase) RemoveFromWishlist(ctx context.Context, userID, productID string) error {
	return u.cartRepo.RemoveFromWishlist(ctx, userID, productID)
}
