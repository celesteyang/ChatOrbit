package main

import (
	"os"
	"strings"
)

const defaultAllowedOrigins = "http://localhost:3000,http://localhost:8080"

// getEnvOrDefault returns the environment variable for the given key or a default value if unset.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseAllowedOrigins returns allowed CORS origins from the ALLOWED_ORIGINS environment variable.
// It falls back to localhost defaults suitable for development.
func parseAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if raw == "" {
		raw = defaultAllowedOrigins
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	return origins
}
