package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type CartWishlistController struct {
	service service.CartWishlistService
}

func NewCartWishlistController(service service.CartWishlistService) *CartWishlistController {
	return &CartWishlistController{service: service}
}

func (cwc *CartWishlistController) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	cart, err := cwc.service.GetCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   fiber.StatusOK,
		"status": "success",
		"cart":   cart,
	})
}

func (cwc *CartWishlistController) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req validation.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	cart, err := cwc.service.AddToCart(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Item added to cart",
		"cart":    cart,
	})
}

func (cwc *CartWishlistController) UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	itemID := c.Params("id")

	var req validation.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	cart, err := cwc.service.UpdateCartItem(userID, itemID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Cart item updated",
		"cart":    cart,
	})
}

func (cwc *CartWishlistController) RemoveCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	itemID := c.Params("id")

	cart, err := cwc.service.RemoveCartItem(userID, itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": "Item removed from cart",
		"cart":    cart,
	})
}

func (cwc *CartWishlistController) GetWishlist(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	items, err := cwc.service.GetWishlist(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     fiber.StatusOK,
		"status":   "success",
		"wishlist": items,
	})
}

func (cwc *CartWishlistController) ToggleWishlist(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("productId")

	added, err := cwc.service.ToggleWishlist(userID, productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	msg := "Removed from wishlist"
	if added {
		msg = "Added to wishlist"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"status":  "success",
		"message": msg,
		"added":   added,
	})
}
