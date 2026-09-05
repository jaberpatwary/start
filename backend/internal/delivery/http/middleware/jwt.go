package middleware

import (
	"net/http"
	"strings"

	"github.com/jaberpatwary/startech/internal/domain"
	"github.com/jaberpatwary/startech/internal/infrastructure/security"
	"github.com/labstack/echo/v4"
)

func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := ""
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else if cookie, err := c.Cookie("auth_token"); err == nil {
				tokenString = cookie.Value
			}

			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication token required")
			}

			claims, err := security.ParseToken(tokenString, jwtSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}

func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			if role != domain.RoleAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "Admin privileges required")
			}
			return next(c)
		}
	}
}
