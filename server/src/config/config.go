package config

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

var (
	AppHost         string
	AppPort         int
	AppEnv          string
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	DBSslMode       string
	DatabaseURL     string
	JWTSecret       string
	JWTExpiresHours int
	ClientURL       string
	UploadDir       string
)

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Default fallback values
	viper.SetDefault("APP_HOST", "0.0.0.0")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "root")
	viper.SetDefault("DB_NAME", "start")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:root@localhost:5432/start?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "startech_secret_key_2026")
	viper.SetDefault("JWT_EXPIRES_HOURS", 72)
	viper.SetDefault("CLIENT_URL", "http://localhost:5173")
	viper.SetDefault("UPLOAD_DIR", "./uploads")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found or unreadable, using environment defaults: %v", err)
	}

	AppHost = viper.GetString("APP_HOST")
	AppPort = viper.GetInt("APP_PORT")
	AppEnv = viper.GetString("APP_ENV")
	DBHost = viper.GetString("DB_HOST")
	DBPort = viper.GetInt("DB_PORT")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
	DBSslMode = viper.GetString("DB_SSLMODE")
	DatabaseURL = viper.GetString("DATABASE_URL")
	JWTSecret = viper.GetString("JWT_SECRET")
	JWTExpiresHours = viper.GetInt("JWT_EXPIRES_HOURS")
	ClientURL = viper.GetString("CLIENT_URL")
	UploadDir = viper.GetString("UPLOAD_DIR")
}

func FiberConfig() fiber.Config {
	return fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"code":    code,
				"status":  "error",
				"message": err.Error(),
			})
		},
	}
}
