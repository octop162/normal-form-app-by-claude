// Package main is the entry point for the normal-form-app server.
// This file is part of the normal-form-app project
package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/handler"
	"github.com/octop162/normal-form-app-by-claude/internal/middleware"
	"github.com/octop162/normal-form-app-by-claude/pkg/config"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	readTimeoutSeconds     = 15
	writeTimeoutSeconds    = 15
	idleTimeoutSeconds     = 60
	shutdownTimeoutSeconds = 30
)

// Application holds all application components
type Application struct {
	UserHandler    *handler.UserHandler
	SessionHandler *handler.SessionHandler
	OptionHandler  *handler.OptionHandler
	AddressHandler *handler.AddressHandler
	PlanHandler    *handler.PlanHandler
	HealthHandler  *handler.HealthHandler
	DB             *sql.DB
	Logger         *logger.Logger
	Config         *config.Config
}

func main() {
	// Initialize application with dependency injection
	app, cleanup, err := wireApp()
	if err != nil {
		panic("Failed to initialize application: " + err.Error())
	}
	defer func() {
		if cleanup != nil {
			cleanup()
		}
	}()

	log := app.Logger
	cfg := app.Config

	log.Infof("Starting normal-form-app server in %s mode", cfg.Server.Mode)
	logger.InitDefaultLogger(cfg.Log.Level)

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	r := setupRouter(app)

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
func setupRouter(app *Application) *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(middleware.SimpleLoggerMiddleware(app.Logger))
	r.Use(middleware.ErrorHandlerMiddleware(app.Logger))
	r.Use(middleware.CORSMiddleware())
	
	// Security middleware
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.InputSanitization())
	r.Use(middleware.RateLimit(100, 1*time.Minute)) // 100 requests per minute
	r.Use(middleware.CSRF())

	// Set up 404 and 405 handlers
	r.NoRoute(middleware.NotFoundMiddleware())
	r.NoMethod(middleware.MethodNotAllowedMiddleware())

	// Health check endpoints
	health := r.Group("/health")
	{
		health.GET("", app.HealthHandler.Health)
		health.GET("/live", app.HealthHandler.LivenessProbe)
		health.GET("/ready", app.HealthHandler.ReadinessProbe)
	}

	// API v1 routes
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data": gin.H{
					"message": "pong",
					"service": "normal-form-app",
					"version": "1.0.0",
				},
			})
		})
		
		// CSRF token endpoint - handled by CSRF middleware
		api.GET("/csrf-token", func(c *gin.Context) {
			// This route is handled by the CSRF middleware
		})

		// User endpoints
		users := api.Group("/users")
		{
			users.POST("", app.UserHandler.CreateUser)
			users.POST("/validate", app.UserHandler.ValidateUser)
			users.GET("/:id", app.UserHandler.GetUser)
			users.PUT("/:id", app.UserHandler.UpdateUser)
			users.DELETE("/:id", app.UserHandler.DeleteUser)
		}

		// Session endpoints
		sessions := api.Group("/sessions")
		{
			sessions.POST("", app.SessionHandler.CreateSession)
			sessions.GET("/:id", app.SessionHandler.GetSession)
			sessions.PUT("/:id", app.SessionHandler.UpdateSession)
			sessions.DELETE("/:id", app.SessionHandler.DeleteSession)
		}

		// Option endpoints
		options := api.Group("/options")
		{
			options.GET("", app.OptionHandler.GetOptions)
			options.POST("/check-inventory", app.OptionHandler.CheckInventory)
			options.GET("/:type", app.OptionHandler.GetOption)
		}

		// Address endpoints
		api.GET("/address/search", app.AddressHandler.SearchAddress)
		api.POST("/region/check", app.AddressHandler.CheckRegion)

		// Prefecture endpoints
		prefectures := api.Group("/prefectures")
		{
			prefectures.GET("", app.AddressHandler.GetPrefectures)
			prefectures.GET("/:name", app.AddressHandler.GetPrefecture)
		}

		// Plan endpoints
		plans := api.Group("/plans")
		{
			plans.GET("", app.PlanHandler.GetPlans)
			plans.GET("/:type", app.PlanHandler.GetPlan)
		}
	}

	return r
}
