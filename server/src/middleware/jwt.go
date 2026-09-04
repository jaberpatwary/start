package middleware

import (
	"app/src/response"
	"app/src/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		var tokenStr string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenStr = c.Cookies("auth_token")
		}

		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorDetails{
				Code:    fiber.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized access, token missing",
			})
		}

		claims, err := utils.ParseJWTToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorDetails{
				Code:    fiber.StatusUnauthorized,
				Status:  "error",
				Message: "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok || role != "ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(response.ErrorDetails{
				Code:    fiber.StatusForbidden,
				Status:  "error",
				Message: "Access denied. Admin rights required",
			})
		}
		return c.Next()
	}
}
