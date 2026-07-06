package config

import (
	"errors"
	"os"
)

const defaultHTTPAddress = ":8080"

var (
	// ErrMissingDatabaseURL is returned when CORDELL_DATABASE_URL is not configured.
	ErrMissingDatabaseURL = errors.New("CORDELL_DATABASE_URL is required")
)

// Config contains runtime settings loaded from the environment.
type Config struct {
	HTTPAddress string
	DatabaseURL string
}

// Load reads runtime settings from environment variables and applies safe defaults.
func Load() (Config, error) {
	databaseURL := os.Getenv("CORDELL_DATABASE_URL")
	if databaseURL == "" {
		return Config{}, ErrMissingDatabaseURL
	}

	return Config{
		HTTPAddress: getEnv("CORDELL_HTTP_ADDRESS", defaultHTTPAddress),
		DatabaseURL: databaseURL,
	}, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
