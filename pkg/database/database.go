// Package database provides database connection and management functionality.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	maxOpenConnections        = 25
	maxIdleConnections        = 25
	connectionMaxLifeMinutes  = 5
	healthCheckTimeoutSeconds = 5
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DB represents the database connection
type DB struct {
	*sql.DB
	config *Config
	log    *logger.Logger
}

// NewDB creates a new database connection
func NewDB(config *Config, log *logger.Logger) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connectionMaxLifeMinutes * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if log != nil {
		log.Info("Database connection established successfully")
	}

	return &DB{
		DB:     db,
		config: config,
		log:    log,
	}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.log != nil {
		d.log.Info("Closing database connection")
	}
	return d.DB.Close()
}

// Ping tests the database connection
func (d *DB) Ping() error {
	return d.DB.Ping()
}

// Stats returns database statistics
func (d *DB) Stats() sql.DBStats {
	return d.DB.Stats()
}

// HealthCheck performs a health check on the database
func (d *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeoutSeconds*time.Second)
	defer cancel()

	if err := d.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// GetConfig returns the database configuration (without password)
func (d *DB) GetConfig() Config {
	config := *d.config
	config.Password = "***" // Hide password for security
	return config
}
