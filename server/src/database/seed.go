package database

import (
	"app/src/model"
	"app/src/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedInitialData(db *gorm.DB) {
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 {
		return // Data already exists
	}

	utils.Log.Info("Seeding initial database data...")

	// Seed Admin & Default User
	adminPass, _ := utils.HashPassword("Admin@123456")
	admin := model.User{
		ID:           uuid.New().String(),
		Name:         "StarTech Admin",
		Email:        "admin@startech.com",
		Phone:        "01700000000",
		PasswordHash: adminPass,
		Role:         model.RoleAdmin,
	}
	db.Create(&admin)

	userPass, _ := utils.HashPassword("User@123456")
	user := model.User{
		ID:           uuid.New().String(),
		Name:         "John Doe",
		Email:        "user@gmail.com",
		Phone:        "01800000000",
		PasswordHash: userPass,
		Role:         model.RoleUser,
	}
	db.Create(&user)

	// Seed Categories
	categories := []model.Category{
		{ID: uuid.New().String(), Name: "Laptop & Notebook", Slug: "laptop-notebook", Icon: "laptop"},
		{ID: uuid.New().String(), Name: "Desktop Component", Slug: "desktop-component", Icon: "cpu"},
		{ID: uuid.New().String(), Name: "Monitor", Slug: "monitor", Icon: "monitor"},
		{ID: uuid.New().String(), Name: "Gaming Gear", Slug: "gaming-gear", Icon: "gamepad"},
	}
	for i := range categories {
		db.Create(&categories[i])
	}

	// Seed Brands
	brands := []model.Brand{
		{ID: uuid.New().String(), Name: "ASUS", Slug: "asus"},
		{ID: uuid.New().String(), Name: "MSI", Slug: "msi"},
		{ID: uuid.New().String(), Name: "Gigabyte", Slug: "gigabyte"},
		{ID: uuid.New().String(), Name: "Intel", Slug: "intel"},
		{ID: uuid.New().String(), Name: "AMD", Slug: "amd"},
	}
	for i := range brands {
		db.Create(&brands[i])
	}

	// Seed Sample Product
	sampleProduct := model.Product{
		ID:               uuid.New().String(),
		Name:             "ASUS ROG Strix G16 Gaming Laptop",
		Slug:             "asus-rog-strix-g16-gaming-laptop",
		SKU:              "ROG-G16-2026",
		CategoryID:       categories[0].ID,
		BrandID:          brands[0].ID,
		Price:            185000,
		Stock:            10,
		ShortDescription: "Intel Core i7 14th Gen, RTX 4060 8GB, 16GB DDR5, 1TB NVMe SSD",
		Description:      "Power up your play with the ROG Strix G16 featuring 14th Gen Intel Core Processor and NVIDIA GeForce RTX 40 Series Laptop GPU.",
		Status:           model.ProductActive,
		Featured:         true,
	}
	db.Create(&sampleProduct)

	utils.Log.Info("Initial data seeding completed")
}
