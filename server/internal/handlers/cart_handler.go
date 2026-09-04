package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jaberpatwary/startech/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{DB: db}
}

func (h *CartHandler) getOrCreateCart(userID string) models.Cart {
	var cart models.Cart
	err := h.DB.Preload("Items.Product.Images").Preload("Items.Product.Brand").Preload("Items.Product.Category").
		Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		cart = models.Cart{ID: uuid.NewString(), UserID: userID, UpdatedAt: time.Now()}
		h.DB.Create(&cart)
	}
	return cart
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userID := c.Get("user_id").(string)
	cart := h.getOrCreateCart(userID)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) AddItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind(&input); err != nil || input.Quantity < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	var product models.Product
	if err := h.DB.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	cart := h.getOrCreateCart(userID)
	var existing models.CartItem
	if err := h.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&existing).Error; err == nil {
		existing.Quantity += input.Quantity
		h.DB.Save(&existing)
	} else {
		item := models.CartItem{ID: uuid.NewString(), CartID: cart.ID, ProductID: input.ProductID, Quantity: input.Quantity}
		h.DB.Create(&item)
	}

	h.DB.Preload("Items.Product.Images").Preload("Items.Product.Brand").Preload("Items.Product.Category").
		Where("user_id = ?", userID).First(&cart)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Item added to cart", "cart": cart})
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	itemID := c.Param("id")
	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := c.Bind(&input); err != nil || input.Quantity < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity")
	}

	cart := h.getOrCreateCart(userID)
	var item models.CartItem
	if err := h.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).First(&item).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Cart item not found")
	}
	item.Quantity = input.Quantity
	h.DB.Save(&item)

	h.DB.Preload("Items.Product.Images").Preload("Items.Product.Brand").Preload("Items.Product.Category").
		Where("user_id = ?", userID).First(&cart)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	itemID := c.Param("id")
	cart := h.getOrCreateCart(userID)
	h.DB.Where("id = ? AND cart_id = ?", itemID, cart.ID).Delete(&models.CartItem{})
	h.DB.Preload("Items.Product.Images").Preload("Items.Product.Brand").Preload("Items.Product.Category").
		Where("user_id = ?", userID).First(&cart)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) ClearCart(userID string) {
	var cart models.Cart
	if err := h.DB.Where("user_id = ?", userID).First(&cart).Error; err == nil {
		h.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	}
}

// Wishlist handler
type WishlistHandler struct {
	DB *gorm.DB
}

func NewWishlistHandler(db *gorm.DB) *WishlistHandler {
	return &WishlistHandler{DB: db}
}

func (h *WishlistHandler) GetWishlist(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var items []models.Wishlist
	h.DB.Preload("Product.Images").Preload("Product.Brand").Preload("Product.Category").
		Where("user_id = ?", userID).Find(&items)
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "wishlist": items})
}

func (h *WishlistHandler) Toggle(c echo.Context) error {
	userID := c.Get("user_id").(string)
	productID := c.Param("productId")

	var existing models.Wishlist
	if err := h.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&existing).Error; err == nil {
		h.DB.Delete(&existing)
		return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Removed from wishlist", "added": false})
	}

	item := models.Wishlist{ID: uuid.NewString(), UserID: userID, ProductID: productID}
	h.DB.Create(&item)
	return c.JSON(http.StatusCreated, echo.Map{"code": 201, "status": "success", "message": "Added to wishlist", "added": true})
}
