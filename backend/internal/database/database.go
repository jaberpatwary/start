package database

import (
	"log"

	"github.com/jaberpatwary/startech/internal/config"
	"github.com/jaberpatwary/startech/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	// GORM log level: Warn in production to prevent leaking sensitive queries and high CPU, Info in development
	logLevel := logger.Info
	if cfg.IsProduction() {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Production database connection pooling
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	log.Printf("Database connection pool configured (MaxOpen: %d, MaxIdle: %d, MaxLifetime: %v)",
		cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)

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

	log.Println("Database migration completed successfully")

	// Seed initial data if database is empty
	if err := Seed(db); err != nil {
		log.Printf("Notice: seeding completed or skipped: %v", err)
	}

	return db, nil
}
