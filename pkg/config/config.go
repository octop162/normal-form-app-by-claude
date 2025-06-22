// Package config provides configuration management functionality.
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/octop162/normal-form-app-by-claude/pkg/database"
)

const (
	defaultPostgresPort = 5432
)

// Config holds all configuration for the application
type Config struct {
	Server      ServerConfig      `json:"server"`
	Database    database.Config   `json:"database"`
	Log         LogConfig         `json:"log"`
	ExternalAPI ExternalAPIConfig `json:"external_api"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
	Mode string `json:"mode"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string `json:"level"`
}

// ExternalAPIConfig holds external API configuration
type ExternalAPIConfig struct {
	InventoryAPI APIConfig `json:"inventory_api"`
	RegionAPI    APIConfig `json:"region_api"`
	AddressAPI   APIConfig `json:"address_api"`
}

// APIConfig holds configuration for a single external API
type APIConfig struct {
	BaseURL    string        `json:"base_url"`
	Timeout    time.Duration `json:"timeout"`
	MaxRetries int           `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load() // .env file not found is not an error

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
			Mode: getEnv("GO_ENV", "development"),
		},
		Database: database.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", defaultPostgresPort),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "normal_form_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		ExternalAPI: ExternalAPIConfig{
			InventoryAPI: APIConfig{
				BaseURL:    getEnv("INVENTORY_API_URL", ""),
				Timeout:    getEnvAsDuration("INVENTORY_API_TIMEOUT", 30*time.Second),
				MaxRetries: getEnvAsInt("INVENTORY_API_MAX_RETRIES", 3),
				RetryDelay: getEnvAsDuration("INVENTORY_API_RETRY_DELAY", 1*time.Second),
			},
			RegionAPI: APIConfig{
				BaseURL:    getEnv("REGION_API_URL", ""),
				Timeout:    getEnvAsDuration("REGION_API_TIMEOUT", 30*time.Second),
				MaxRetries: getEnvAsInt("REGION_API_MAX_RETRIES", 3),
				RetryDelay: getEnvAsDuration("REGION_API_RETRY_DELAY", 1*time.Second),
			},
			AddressAPI: APIConfig{
				BaseURL:    getEnv("ADDRESS_API_URL", ""),
				Timeout:    getEnvAsDuration("ADDRESS_API_TIMEOUT", 30*time.Second),
				MaxRetries: getEnvAsInt("ADDRESS_API_MAX_RETRIES", 3),
				RetryDelay: getEnvAsDuration("ADDRESS_API_RETRY_DELAY", 1*time.Second),
			},
		},
	}

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDuration gets an environment variable as duration or returns a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Mode == "production"
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Mode == "development"
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}
