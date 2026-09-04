package database

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, config.DBSslMode,
	)

	// Fallback to DatabaseURL if specified
	if config.DatabaseURL != "" {
		dsn = config.DatabaseURL
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		utils.Log.Fatalf("Failed to connect to database: %v", err)
	}

	utils.Log.Info("Database connection established successfully")

	// Run Auto Migrations
	err = db.AutoMigrate(
		&model.User{},
		&model.Address{},
		&model.Category{},
		&model.Brand{},
		&model.Product{},
		&model.ProductImage{},
		&model.Cart{},
		&model.CartItem{},
		&model.Wishlist{},
		&model.Order{},
		&model.OrderItem{},
		&model.Review{},
		&model.Coupon{},
		&model.Banner{},
		&model.Setting{},
	)

	if err != nil {
		utils.Log.Fatalf("Failed to run database migrations: %v", err)
	}

	utils.Log.Info("Database auto-migration completed successfully")

	// Seed initial data
	SeedInitialData(db)

	return db
}
