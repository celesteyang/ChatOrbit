package main

import (
	"os"
	"strings"
)

// getEnvOrDefault returns the environment variable for the given key or a default value if unset.
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseBoolEnv parses a boolean environment variable, falling back to a default when empty.
func parseBoolEnv(key, defaultValue string) bool {
	val := getEnvOrDefault(key, defaultValue)
	switch val {
	case "1", "true", "TRUE", "True", "yes", "YES", "Yes":
		return true
	}
	return false
}

// parseAllowedOrigins returns a list of allowed origins from the ALLOWED_ORIGINS environment variable.
// The variable should be a comma-separated list like "https://app.example.com,http://localhost:3000".
// If unset, it falls back to a local development default.
func parseAllowedOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
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
