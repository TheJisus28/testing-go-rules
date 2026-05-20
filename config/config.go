// Package config loads environment variables used at application startup.
package config

import (
	"os"

	"github.com/donca/user-crud/pkg/kit/enums"
	"github.com/joho/godotenv"
)

// LoadEnv reads .env when present and applies safe defaults for local development.
func LoadEnv() {
	_ = godotenv.Load()
	setDefault(enums.AppPort, "8080")
	setDefault(enums.LogLevel, "info")
	setDefault(enums.DatabaseURL, "postgres://social:social@localhost:5432/socialnet?sslmode=disable")
	setDefault(enums.JWTSecret, "dev-secret-change-in-production")
}

func setDefault(key, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
