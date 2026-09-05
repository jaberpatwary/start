package main

import (
	"log"

	"github.com/jaberpatwary/startech/internal/config"
	"github.com/jaberpatwary/startech/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env", "../.env", "../../.env")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seed finished successfully.")
}
