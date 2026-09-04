package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	orderService service.OrderService
}

func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (oc *OrderController) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req validation.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	order, err := oc.orderService.CreateOrder(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"status":  "success",
		"message": "Order placed successfully",
		"order":   order,
	})
}

func (oc *OrderController) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orders, err := oc.orderService.GetUserOrders(userID)
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
		"orders": orders,
	})
}

func (oc *OrderController) GetOrderByID(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	orderID := c.Params("id")

	order, err := oc.orderService.GetOrderByID(userID, role, orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
			Code:    fiber.StatusNotFound,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":   fiber.StatusOK,
		"status": "success",
		"order":  order,
	})
}

func (oc *OrderController) AdminGetAllOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	status := c.Query("status", "")

	res, err := oc.orderService.AdminGetAllOrders(page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorDetails{
			Code:    fiber.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (oc *OrderController) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")
	var req validation.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "error",
			Message: "Invalid request payload",
		})
	}

	order, err := oc.orderService.UpdateOrderStatus(orderID, &req)
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
		"message": "Order status updated successfully",
		"order":   order,
	})
}
