package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	ShutdownTO  time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://product:product@localhost:5433/product?sslmode=disable"),
		ShutdownTO:  10 * time.Second,
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
