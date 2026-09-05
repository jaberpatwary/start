package database

import (
	"log"

	"github.com/jaberpatwary/startech/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	err = db.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Category{},
		&models.Brand{},
		&models.Product{},
		&models.ProductImage{},
		&models.Cart{},
		&models.CartItem{},
		&models.Wishlist{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.Coupon{},
		&models.Banner{},
		&models.Setting{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database migration completed")

	// Seed initial data
	if err := Seed(db); err != nil {
		log.Printf("Warning: seeding failed: %v", err)
	}

	return db, nil
}
