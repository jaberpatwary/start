package utils

import (
	"app/src/response"

	"github.com/gofiber/fiber/v2"
)

func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(response.ErrorDetails{
		Code:    fiber.StatusNotFound,
		Status:  "error",
		Message: "Requested route not found",
	})
}
