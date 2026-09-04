package routes

import (
	"github.com/jaberpatwary/startech/internal/handlers"
	mw "github.com/jaberpatwary/startech/internal/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, db *gorm.DB, uploadDir string) {
	// Init handlers
	authH := handlers.NewAuthHandler(db)
	productH := handlers.NewProductHandler(db)
	catalogH := handlers.NewCatalogHandler(db)
	cartH := handlers.NewCartHandler(db)
	wishlistH := handlers.NewWishlistHandler(db)
	orderH := handlers.NewOrderHandler(db, cartH)
	adminH := handlers.NewAdminHandler(db)
	uploadH := handlers.NewUploadHandler(uploadDir)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, echo.Map{"status": "ok", "message": "StarTech API is running"})
	})

	v1 := e.Group("/api/v1")

	// --- Public Routes ---
	// Auth
	v1.POST("/auth/register", authH.Register)
	v1.POST("/auth/login", authH.Login)
	v1.POST("/auth/logout", authH.Logout)

	// Products (public)
	v1.GET("/products", productH.GetAll)
	v1.GET("/products/:id", productH.GetByID)
	v1.GET("/products/slug/:slug", productH.GetBySlug)

	// Categories & Brands (public)
	v1.GET("/categories", catalogH.GetCategories)
	v1.GET("/brands", catalogH.GetBrands)

	// Banners (public)
	v1.GET("/banners", catalogH.GetBanners)

	// Coupon validation (public)
	v1.GET("/coupons/validate", catalogH.ValidateCoupon)

	// --- Protected User Routes ---
	user := v1.Group("", mw.JWTAuth())

	// Auth (protected)
	user.GET("/auth/me", authH.GetMe)
	user.PUT("/users/profile", authH.UpdateProfile)
	user.POST("/users/addresses", authH.AddAddress)
	user.DELETE("/users/addresses/:id", authH.DeleteAddress)

	// Cart
	user.GET("/cart", cartH.GetCart)
	user.POST("/cart/items", cartH.AddItem)
	user.PUT("/cart/items/:id", cartH.UpdateItem)
	user.DELETE("/cart/items/:id", cartH.RemoveItem)

	// Wishlist
	user.GET("/wishlist", wishlistH.GetWishlist)
	user.POST("/wishlist/:productId", wishlistH.Toggle)

	// Orders
	user.POST("/orders", orderH.CreateOrder)
	user.GET("/orders", orderH.GetMyOrders)
	user.GET("/orders/:id", orderH.GetOrderByID)

	// Reviews (submit)
	user.POST("/products/:id/reviews", productH.AddReview)

	// --- Admin Routes ---
	admin := v1.Group("/admin", mw.JWTAuth(), mw.AdminOnly())

	// Dashboard
	admin.GET("/dashboard", adminH.GetDashboardStats)

	// Admin: Products
	admin.POST("/products", productH.Create)
	admin.PUT("/products/:id", productH.Update)
	admin.DELETE("/products/:id", productH.Delete)

	// Admin: Categories
	admin.POST("/categories", catalogH.CreateCategory)
	admin.PUT("/categories/:id", catalogH.UpdateCategory)
	admin.DELETE("/categories/:id", catalogH.DeleteCategory)

	// Admin: Brands
	admin.POST("/brands", catalogH.CreateBrand)
	admin.PUT("/brands/:id", catalogH.UpdateBrand)
	admin.DELETE("/brands/:id", catalogH.DeleteBrand)

	// Admin: Banners
	admin.GET("/banners", catalogH.AdminGetBanners)
	admin.POST("/banners", catalogH.CreateBanner)
	admin.PUT("/banners/:id", catalogH.UpdateBanner)
	admin.DELETE("/banners/:id", catalogH.DeleteBanner)

	// Admin: Coupons
	admin.GET("/coupons", catalogH.GetCoupons)
	admin.POST("/coupons", catalogH.CreateCoupon)
	admin.PUT("/coupons/:id", catalogH.UpdateCoupon)
	admin.DELETE("/coupons/:id", catalogH.DeleteCoupon)

	// Admin: Orders
	admin.GET("/orders", orderH.AdminGetAllOrders)
	admin.GET("/orders/:id", orderH.GetOrderByID)
	admin.PATCH("/orders/:id/status", orderH.UpdateOrderStatus)

	// Admin: Users
	admin.GET("/users", adminH.GetAllUsers)
	admin.PATCH("/users/:id/toggle-block", adminH.ToggleBlockUser)

	// Admin: Reviews
	admin.GET("/reviews", adminH.GetAllReviews)
	admin.PATCH("/reviews/:id/status", adminH.UpdateReviewStatus)

	// Admin: Inventory
	admin.GET("/inventory", adminH.GetInventory)

	// Upload (admin only)
	admin.POST("/upload", uploadH.UploadImage)
}
