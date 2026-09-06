package http

import (
	"net/http"
	"strings"

	handler "github.com/jaberpatwary/startech/internal/delivery/http/handler"
	mw "github.com/jaberpatwary/startech/internal/delivery/http/middleware"
	"github.com/jaberpatwary/startech/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Router struct {
	authHandler    *handler.AuthHandler
	productHandler *handler.ProductHandler
	catalogHandler *handler.CatalogHandler
	cartHandler    *handler.CartHandler
	orderHandler   *handler.OrderHandler
	adminHandler   *handler.AdminHandler
	jwtSecret      string
}

func NewRouter(
	authUC usecase.AuthUsecase,
	productUC usecase.ProductUsecase,
	catalogUC usecase.CatalogUsecase,
	cartUC usecase.CartUsecase,
	orderUC usecase.OrderUsecase,
	adminUC usecase.AdminUsecase,
	uploadDir string,
	jwtSecret string,
) *Router {
	return &Router{
		authHandler:    handler.NewAuthHandler(authUC),
		productHandler: handler.NewProductHandler(productUC),
		catalogHandler: handler.NewCatalogHandler(catalogUC),
		cartHandler:    handler.NewCartHandler(cartUC),
		orderHandler:   handler.NewOrderHandler(orderUC),
		adminHandler:   handler.NewAdminHandler(adminUC, uploadDir),
		jwtSecret:      jwtSecret,
	}
}

func (r *Router) Register(e *echo.Echo) {
	jwtMW := mw.JWTAuth(r.jwtSecret)
	adminMW := mw.RequireAdmin()

	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		accept := c.Request().Header.Get("Accept")
		if strings.Contains(accept, "text/html") {
			return c.HTML(http.StatusOK, `<!DOCTYPE html><html><head><meta charset="utf-8"><title>MI-Tech API</title><style>body{font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh;margin:0;background:#0f172a;color:#fff;text-align:center}a{color:#00C2FF;text-decoration:none;font-weight:bold}a:hover{text-decoration:underline}.card{background:#1e293b;padding:32px;border-radius:12px;border:1px solid #334155;max-width:480px}</style></head><body><div class="card"><h1>🚀 MI-Tech API Server</h1><p style="color:#94a3b8">Backend server is running on port 8090 with Clean Architecture.</p><p style="margin-top:24px"><a href="http://localhost:8888" style="background:#ef4a23;color:#fff;padding:12px 24px;border-radius:6px;display:inline-block">Open Frontend (http://localhost:8888) →</a></p></div></body></html>`)
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":   "ok",
			"message":  "MI-Tech API server is running",
			"frontend": "http://localhost:8888",
			"health":   "/health",
		})
	})

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"status": "ok", "message": "MI-Tech API is running"})
	})

	v1 := e.Group("/api/v1")

	// --- Public Routes ---
	authLimiter := mw.AuthRateLimit()
	v1.POST("/auth/register", r.authHandler.Register, authLimiter)
	v1.POST("/auth/login", r.authHandler.Login, authLimiter)
	v1.POST("/auth/logout", r.authHandler.Logout)

	// Products (public)
	v1.GET("/products", r.productHandler.GetAll)
	v1.GET("/products/slug/:slug", r.productHandler.GetBySlug)
	v1.GET("/products/:id", r.productHandler.GetByID)

	// Categories & Brands (public)
	v1.GET("/categories", r.catalogHandler.GetCategories)
	v1.GET("/brands", r.catalogHandler.GetBrands)

	// Banners (public)
	v1.GET("/banners", r.catalogHandler.GetBanners)

	// Coupon validation (public)
	v1.GET("/coupons/validate", r.catalogHandler.ValidateCoupon)

	// Order tracking (public)
	v1.GET("/orders/track/:number", r.orderHandler.TrackOrderByNumber)

	// Public store settings
	v1.GET("/settings", r.adminHandler.GetSettings)

	// --- Protected User Routes ---
	user := v1.Group("", jwtMW)

	// Auth (protected)
	user.GET("/auth/me", r.authHandler.GetMe)
	user.PUT("/users/profile", r.authHandler.UpdateProfile)
	user.POST("/users/addresses", r.authHandler.AddAddress)
	user.DELETE("/users/addresses/:id", r.authHandler.DeleteAddress)

	// Cart
	user.GET("/cart", r.cartHandler.GetCart)
	user.POST("/cart/items", r.cartHandler.AddItem)
	user.PUT("/cart/items/:id", r.cartHandler.UpdateItem)
	user.DELETE("/cart/items/:id", r.cartHandler.RemoveItem)

	// Wishlist
	user.GET("/wishlist", r.cartHandler.GetWishlist)
	user.POST("/wishlist/:productId", r.cartHandler.ToggleWishlist)

	// Orders
	user.POST("/orders", r.orderHandler.CreateOrder)
	user.GET("/orders", r.orderHandler.GetMyOrders)
	user.GET("/orders/:id", r.orderHandler.GetOrderByID)

	// Reviews (submit)
	user.POST("/products/:id/reviews", r.productHandler.AddReview)

	// --- Admin Routes ---
	admin := v1.Group("/admin", jwtMW, adminMW)

	// Dashboard
	admin.GET("/dashboard", r.adminHandler.GetDashboardStats)

	// Admin: Products
	admin.POST("/products", r.productHandler.Create)
	admin.PUT("/products/:id", r.productHandler.Update)
	admin.DELETE("/products/:id", r.productHandler.Delete)

	// Admin: Categories
	admin.POST("/categories", r.catalogHandler.CreateCategory)
	admin.PUT("/categories/:id", r.catalogHandler.UpdateCategory)
	admin.DELETE("/categories/:id", r.catalogHandler.DeleteCategory)

	// Admin: Brands
	admin.POST("/brands", r.catalogHandler.CreateBrand)
	admin.PUT("/brands/:id", r.catalogHandler.UpdateBrand)
	admin.DELETE("/brands/:id", r.catalogHandler.DeleteBrand)

	// Admin: Banners
	admin.GET("/banners", r.catalogHandler.GetBanners)
	admin.POST("/banners", r.catalogHandler.CreateBanner)
	admin.DELETE("/banners/:id", r.catalogHandler.DeleteBanner)

	// Admin: Coupons
	admin.POST("/coupons", r.catalogHandler.CreateCoupon)
	admin.DELETE("/coupons/:id", r.catalogHandler.DeleteCoupon)

	// Admin: Orders
	admin.GET("/orders", r.orderHandler.AdminGetAllOrders)
	admin.GET("/orders/:id", r.orderHandler.GetOrderByID)
	admin.PATCH("/orders/:id/status", r.orderHandler.UpdateOrderStatus)
	admin.PATCH("/orders/:id/payment", r.orderHandler.UpdatePaymentStatus)

	// Admin: Users
	admin.GET("/users", r.adminHandler.GetAllUsers)
	admin.PATCH("/users/:id/toggle-block", r.adminHandler.ToggleBlockUser)

	// Admin: Reviews
	admin.GET("/reviews", r.adminHandler.GetAllReviews)
	admin.PATCH("/reviews/:id/status", r.adminHandler.UpdateReviewStatus)

	// Admin: Inventory & Reports
	admin.GET("/inventory", r.adminHandler.GetInventory)
	admin.GET("/reports", r.adminHandler.GetReports)

	// Admin: Settings
	admin.GET("/settings", r.adminHandler.GetSettings)
	admin.PUT("/settings", r.adminHandler.UpdateSettings)

	// Upload (admin only)
	admin.POST("/upload", r.adminHandler.UploadImage)
}
