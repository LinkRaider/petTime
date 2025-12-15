package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	JWT         JWTConfig
}

type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Env:         getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://pettime:pettime@localhost:5432/pettime?sslmode=disable"),
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "dev-secret-change-in-production"),
			AccessTokenTTL:   getDurationEnv("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL:  getDurationEnv("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if minutes, err := strconv.Atoi(value); err == nil {
			return time.Duration(minutes) * time.Minute
		}
	}
	return defaultValue
}
