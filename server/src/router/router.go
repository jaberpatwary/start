package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	// Initialize Services
	userService := service.NewUserService(db)
	categoryBrandService := service.NewCategoryBrandService(db)
	productService := service.NewProductService(db)
	cartWishlistService := service.NewCartWishlistService(db)
	orderService := service.NewOrderService(db, cartWishlistService)

	// Initialize Controllers
	authController := controller.NewAuthController(userService)
	categoryBrandController := controller.NewCategoryBrandController(categoryBrandService)
	productController := controller.NewProductController(productService)
	cartWishlistController := controller.NewCartWishlistController(cartWishlistService)
	orderController := controller.NewOrderController(orderService)

	// API Group /v1
	v1 := app.Group("/v1")

	// Health Check
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "StarTech Backend API is running smoothly",
		})
	})

	// Public Auth Endpoints
	auth := v1.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)

	// Public Catalog Endpoints
	v1.Get("/categories", categoryBrandController.GetCategories)
	v1.Get("/brands", categoryBrandController.GetBrands)

	products := v1.Group("/products")
	products.Get("/", productController.GetAll)
	products.Get("/slug/:slug", productController.GetBySlug)
	products.Get("/:id", productController.GetByID)

	// Protected Endpoints (Requires User JWT)
	protected := v1.Group("", middleware.JWTMiddleware())

	// User Profile
	protected.Get("/auth/me", authController.GetMe)
	protected.Put("/users/profile", authController.UpdateProfile)

	// Cart Operations
	cart := protected.Group("/cart")
	cart.Get("/", cartWishlistController.GetCart)
	cart.Post("/items", cartWishlistController.AddToCart)
	cart.Put("/items/:id", cartWishlistController.UpdateCartItem)
	cart.Delete("/items/:id", cartWishlistController.RemoveCartItem)

	// Wishlist Operations
	wishlist := protected.Group("/wishlist")
	wishlist.Get("/", cartWishlistController.GetWishlist)
	wishlist.Post("/:productId", cartWishlistController.ToggleWishlist)

	// Order Operations
	orders := protected.Group("/orders")
	orders.Post("/", orderController.CreateOrder)
	orders.Get("/", orderController.GetUserOrders)
	orders.Get("/:id", orderController.GetOrderByID)

	// Admin Endpoints (Requires Admin Role)
	admin := protected.Group("/admin", middleware.AdminOnly())
	admin.Get("/users", authController.GetAllUsers)

	// Admin Catalog Operations
	v1.Post("/categories", middleware.JWTMiddleware(), middleware.AdminOnly(), categoryBrandController.CreateCategory)
	v1.Post("/brands", middleware.JWTMiddleware(), middleware.AdminOnly(), categoryBrandController.CreateBrand)

	productsAdmin := v1.Group("/products", middleware.JWTMiddleware(), middleware.AdminOnly())
	productsAdmin.Post("/", productController.CreateProduct)
	productsAdmin.Put("/:id", productController.UpdateProduct)
	productsAdmin.Delete("/:id", productController.DeleteProduct)

	// Admin Orders
	adminOrders := admin.Group("/orders")
	adminOrders.Get("/", orderController.AdminGetAllOrders)
	adminOrders.Patch("/:id/status", orderController.UpdateOrderStatus)
}
