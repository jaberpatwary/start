package handler

import (
	"net/http"

	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartUsecase usecase.CartUsecase
}

func NewCartHandler(cartUsecase usecase.CartUsecase) *CartHandler {
	return &CartHandler{cartUsecase: cartUsecase}
}

func (h *CartHandler) GetCart(c echo.Context) error {
	userID := c.Get("user_id").(string)
	cart, err := h.cartUsecase.GetCart(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get cart")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) AddItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	var req struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	cart, err := h.cartUsecase.AddToCart(c.Request().Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) UpdateItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	itemID := c.Param("id")
	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	cart, err := h.cartUsecase.UpdateItem(c.Request().Context(), userID, itemID, req.Quantity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) RemoveItem(c echo.Context) error {
	userID := c.Get("user_id").(string)
	itemID := c.Param("id")
	cart, err := h.cartUsecase.RemoveItem(c.Request().Context(), userID, itemID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "cart": cart})
}

func (h *CartHandler) GetWishlist(c echo.Context) error {
	userID := c.Get("user_id").(string)
	wishlist, err := h.cartUsecase.GetWishlist(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get wishlist")
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "wishlist": wishlist})
}

func (h *CartHandler) ToggleWishlist(c echo.Context) error {
	userID := c.Get("user_id").(string)
	productID := c.Param("productId")

	// Try remove first; if not found, add
	err := h.cartUsecase.RemoveFromWishlist(c.Request().Context(), userID, productID)
	if err != nil {
		// Not in wishlist, add it
		if addErr := h.cartUsecase.AddToWishlist(c.Request().Context(), userID, productID); addErr != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, addErr.Error())
		}
		return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Added to wishlist", "action": "added"})
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Removed from wishlist", "action": "removed"})
}

// ClearCart is a helper used internally by order creation; exposed here for completeness
func (h *CartHandler) ClearCart(c echo.Context) error {
	userID := c.Get("user_id").(string)
	if err := h.cartUsecase.ClearCart(c.Request().Context(), userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{"code": 200, "status": "success", "message": "Cart cleared"})
}
