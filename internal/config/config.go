package config

import "os"

type Config struct {
	AppEnv      string
	Port        string
	LogLevel    string
	DatabaseURL string
}

func Load() Config {
	return Config{
		AppEnv:      get("APP_ENV", "development"),
		Port:        get("APP_PORT", "8080"),
		LogLevel:    get("LOG_LEVEL", "info"),
		DatabaseURL: get("DATABASE_URL", "postgres://app:app@localhost:5433/app?sslmode=disable"),
	}
}

func get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
