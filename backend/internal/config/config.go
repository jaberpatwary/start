package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv            string
	Port              string
	DatabaseURL       string
	JWTSecret         string
	ClientURL         string
	UploadDir         string
	AllowedOrigins    []string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
}

func (c *Config) IsProduction() bool {
	env := strings.ToLower(strings.TrimSpace(c.AppEnv))
	return env == "production" || env == "prod"
}

func Load() (*Config, error) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:root@localhost:5432/start?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "startech_secret_key_2026"
	}

	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:8888"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Parse Allowed Origins for CORS
	allowedOriginsRaw := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsRaw != "" {
		parts := strings.Split(allowedOriginsRaw, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
	}
	// Always include clientURL
	if clientURL != "" {
		allowedOrigins = append(allowedOrigins, clientURL)
	}
	// Add common local origins in non-production
	if appEnv != "production" {
		allowedOrigins = append(allowedOrigins,
			"http://localhost:5173",
			"http://localhost:3000",
			"http://localhost:8888",
			"http://localhost:8090",
			"http://127.0.0.1:8888",
			"http://127.0.0.1:8090",
		)
	}

	// Database Connection Pooling Defaults
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	connMaxLifetime := time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 5)) * time.Minute
	connMaxIdleTime := time.Duration(getEnvAsInt("DB_CONN_MAX_IDLE_MIN", 15)) * time.Minute

	return &Config{
		AppEnv:            appEnv,
		Port:              port,
		DatabaseURL:       dbURL,
		JWTSecret:         jwtSecret,
		ClientURL:         clientURL,
		UploadDir:         uploadDir,
		AllowedOrigins:    uniqueStrings(allowedOrigins),
		DBMaxOpenConns:    maxOpenConns,
		DBMaxIdleConns:    maxIdleConns,
		DBConnMaxLifetime: connMaxLifetime,
		DBConnMaxIdleTime: connMaxIdleTime,
	}, nil
}

func getEnvAsInt(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

func uniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range slice {
		if _, value := keys[entry]; !value && entry != "" {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
