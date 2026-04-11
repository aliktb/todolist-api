package config

import (
	"os"
	"time"
)

type Config struct {
	Address          string
	DatabaseURL      string
	AllowedOrigin    string
	ReadinessTimeout time.Duration
}

func Load() Config {
	return Config{
		Address:          getEnv("APP_ADDR", ":8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todolist?sslmode=disable"),
		AllowedOrigin:    getEnv("ALLOWED_ORIGIN", "*"),
		ReadinessTimeout: getDurationEnv("READINESS_TIMEOUT", 2*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return fallback
}
