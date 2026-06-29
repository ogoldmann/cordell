package config

import "os"

const defaultHTTPAddress = ":8080"

// Config contains runtime settings loaded from the environment.
type Config struct {
	HTTPAddress string
}

// Load reads runtime settings from environment variables and applies safe defaults.
func Load() Config {
	return Config{
		HTTPAddress: getEnv("CORDELL_HTTP_ADDRESS", defaultHTTPAddress),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
