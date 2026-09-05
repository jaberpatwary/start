package main

import (
	"log"
	"os"

	"github.com/jaberpatwary/startech/internal/config"
	"github.com/jaberpatwary/startech/internal/database"
	deliveryhttp "github.com/jaberpatwary/startech/internal/delivery/http"
	pgRepo "github.com/jaberpatwary/startech/internal/repository/postgres"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env from current dir or parent
	if err := godotenv.Load(); err != nil {
		_ = godotenv.Load("../.env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect DB
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create upload directory
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal("Failed to create upload dir:", err)
	}

	// ─── Dependency Injection (Composition Root) ───────────────────────────────
	// Repositories (Infrastructure Layer)
	userRepo    := pgRepo.NewUserRepository(db)
	productRepo := pgRepo.NewProductRepository(db)
	catalogRepo := pgRepo.NewCatalogRepository(db)
	cartRepo    := pgRepo.NewCartRepository(db)
	orderRepo   := pgRepo.NewOrderRepository(db)
	adminRepo   := pgRepo.NewAdminRepository(db)

	// Use Cases (Application Layer)
	authUC    := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	productUC := usecase.NewProductUsecase(productRepo)
	catalogUC := usecase.NewCatalogUsecase(catalogRepo)
	cartUC    := usecase.NewCartUsecase(cartRepo)
	orderUC   := usecase.NewOrderUsecase(orderRepo, cartRepo, productRepo, catalogRepo)
	adminUC   := usecase.NewAdminUsecase(adminRepo, userRepo)
	// ──────────────────────────────────────────────────────────────────────────

	// Echo setup
	e := echo.New()
	e.HideBanner = true

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.ClientURL, "http://localhost:5173", "http://localhost:3000", "http://localhost:8888", "http://localhost:8090"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Serve uploaded files
	e.Static("/uploads", cfg.UploadDir)

	// Register routes via Clean Architecture Router
	router := deliveryhttp.NewRouter(authUC, productUC, catalogUC, cartUC, orderUC, adminUC, cfg.UploadDir, cfg.JWTSecret)
	router.Register(e)

	log.Printf("MI-Tech API server listening on :%s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
