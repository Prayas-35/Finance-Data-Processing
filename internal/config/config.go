package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds runtime settings loaded from environment variables.
type Config struct {
	AppPort           string
	DatabaseURL       string
	JWTSecret         string
	JWTAccessTokenTTL time.Duration
	DefaultPageLimit  int
	MaxPageLimit      int
}

// Load reads configuration from environment.
func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	cfg := Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTAccessTokenTTL: getDuration("JWT_ACCESS_TOKEN_TTL", time.Hour),
		DefaultPageLimit:  getInt("DEFAULT_PAGE_LIMIT", 25),
		MaxPageLimit:      getInt("MAX_PAGE_LIMIT", 100),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.DefaultPageLimit <= 0 {
		cfg.DefaultPageLimit = 25
	}

	if cfg.MaxPageLimit < cfg.DefaultPageLimit {
		cfg.MaxPageLimit = 100
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return parsed
}
