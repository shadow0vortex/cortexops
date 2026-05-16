package config

import (
	"fmt"
	"os"
	"strconv"
)

// AppConfig represents global shared configuration needed by most services.
type AppConfig struct {
	Environment string
	Port        int
	LogLevel    string
	NatsURL     string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*AppConfig, error) {
	portStr := getEnvOrDefault("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT environment variable: %w", err)
	}

	return &AppConfig{
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Port:        port,
		LogLevel:    getEnvOrDefault("LOG_LEVEL", "INFO"),
		NatsURL:     getEnvOrDefault("NATS_URL", "nats://localhost:4222"),
	}, nil
}

// getEnvOrDefault gets an environment variable or returns the default value.
func getEnvOrDefault(key, defaultValue string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultValue
}
