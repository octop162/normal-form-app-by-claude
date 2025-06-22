// Package main is the entry point for the normal-form-app server.
// This file is part of the normal-form-app project
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/handler"
	"github.com/octop162/normal-form-app-by-claude/internal/middleware"
	"github.com/octop162/normal-form-app-by-claude/pkg/config"
	"github.com/octop162/normal-form-app-by-claude/pkg/database"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	readTimeoutSeconds     = 15
	writeTimeoutSeconds    = 15
	idleTimeoutSeconds     = 60
	shutdownTimeoutSeconds = 30
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Log.Level)
	logger.InitDefaultLogger(cfg.Log.Level)

	log.Infof("Starting normal-form-app server in %s mode", cfg.Server.Mode)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connection
	var db *database.DB
	if cfg.Database.Host != "" {
		db, err = database.NewDB(&cfg.Database, log)
		if err != nil {
			log.WithError(err).Error("Failed to connect to database")
			// Don't exit, allow server to run without database for health checks
		} else {
			defer func() {
				if closeErr := db.Close(); closeErr != nil {
					log.WithError(closeErr).Error("Failed to close database connection")
				}
			}()
		}
	}

	// Create handlers
	healthHandler := handler.NewHealthHandler(db, log)

	// Create router
	r := setupRouter(log, healthHandler)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  readTimeoutSeconds * time.Second,
		WriteTimeout: writeTimeoutSeconds * time.Second,
		IdleTimeout:  idleTimeoutSeconds * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Server starting on %s", cfg.GetServerAddress())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Give outstanding requests a deadline to complete
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeoutSeconds*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Server exited")
}

// setupRouter configures and returns the Gin router
func setupRouter(log *logger.Logger, healthHandler *handler.HealthHandler) *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(middleware.SimpleLoggerMiddleware(log))
	r.Use(middleware.ErrorHandlerMiddleware(log))
	r.Use(middleware.CORSMiddleware())

	// Set up 404 and 405 handlers
	r.NoRoute(middleware.NotFoundMiddleware())
	r.NoMethod(middleware.MethodNotAllowedMiddleware())

	// Health check endpoints
	health := r.Group("/health")
	{
		health.GET("", healthHandler.Health)
		health.GET("/live", healthHandler.LivenessProbe)
		health.GET("/ready", healthHandler.ReadinessProbe)
	}

	// API v1 routes
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"service": "normal-form-app",
				"version": "1.0.0",
			})
		})
	}

	return r
}
