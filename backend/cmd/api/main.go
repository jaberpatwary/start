package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jaberpatwary/startech/internal/config"
	"github.com/jaberpatwary/startech/internal/database"
	deliveryhttp "github.com/jaberpatwary/startech/internal/delivery/http"
	pgRepo "github.com/jaberpatwary/startech/internal/repository/postgres"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	// Load .env from current dir or parent directories
	_ = godotenv.Load(".env", "../.env", "../../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect DB with connection pooling
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create upload directory
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload dir:", err)
	}

	// ─── Dependency Injection (Composition Root) ───────────────────────────────
	// Repositories (Infrastructure Layer)
	userRepo := pgRepo.NewUserRepository(db)
	productRepo := pgRepo.NewProductRepository(db)
	catalogRepo := pgRepo.NewCatalogRepository(db)
	cartRepo := pgRepo.NewCartRepository(db)
	orderRepo := pgRepo.NewOrderRepository(db)
	adminRepo := pgRepo.NewAdminRepository(db)

	// Use Cases (Application Layer)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	productUC := usecase.NewProductUsecase(productRepo)
	catalogUC := usecase.NewCatalogUsecase(catalogRepo)
	cartUC := usecase.NewCartUsecase(cartRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, cartRepo, productRepo, catalogRepo)
	adminUC := usecase.NewAdminUsecase(adminRepo, userRepo)
	// ──────────────────────────────────────────────────────────────────────────

	// Echo setup
	e := echo.New()
	e.HideBanner = true

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Security Headers (Defense in depth)
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
	}))

	// Limit request payload size to prevent DoS via large bodies
	e.Use(middleware.BodyLimit("25M"))

	// Rate Limiting (50 requests/sec with burst buffer)
	rateStore := middleware.NewRateLimiterMemoryStore(rate.Limit(50))
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: func(c echo.Context) bool {
			// Skip health check and static assets from rate limiting
			path := c.Path()
			return strings.HasPrefix(path, "/health") || strings.HasPrefix(path, "/uploads")
		},
		Store: rateStore,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, echo.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please slow down.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, echo.Map{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	}))

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Serve uploaded files
	e.Static("/uploads", cfg.UploadDir)

	// Enhanced Production Health Check with DB connectivity test
	e.GET("/health", func(c echo.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, echo.Map{
				"status":   "unhealthy",
				"database": "error getting db instance",
			})
		}
		if err := sqlDB.Ping(); err != nil {
			return c.JSON(http.StatusServiceUnavailable, echo.Map{
				"status":   "unhealthy",
				"database": "unreachable",
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":      "ok",
			"environment": cfg.AppEnv,
			"database":    "healthy",
			"uptime":      time.Now().Format(time.RFC3339),
		})
	})

	// Register routes via Clean Architecture Router
	router := deliveryhttp.NewRouter(authUC, productUC, catalogUC, cartUC, orderUC, adminUC, cfg.UploadDir, cfg.JWTSecret)
	router.Register(e)

	// ─── Graceful Shutdown ─────────────────────────────────────────────────────
	go func() {
		log.Printf("MI-Tech API server (%s) listening on :%s", cfg.AppEnv, cfg.Port)
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
