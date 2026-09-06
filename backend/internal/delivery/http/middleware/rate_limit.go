package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// AuthRateLimit restricts sensitive auth attempts (login/register) to prevent brute-force attacks
func AuthRateLimit() echo.MiddlewareFunc {
	store := middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:      rate.Every(6 * time.Second), // 10 attempts per minute
		Burst:     10,
		ExpiresIn: 5 * time.Minute,
	})

	return middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Store: store,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, echo.Map{
				"code":    429,
				"status":  "error",
				"message": "Too many requests. Please slow down.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, echo.Map{
				"code":    429,
				"status":  "error",
				"message": "Too many attempts from your IP. Please wait a minute before trying again.",
			})
		},
	})
}
