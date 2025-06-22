// Package handler provides HTTP request handlers.
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/pkg/database"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	statusHealthy       = "healthy"
	statusUnhealthy     = "unhealthy"
	statusNotConfigured = "not configured"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db  *database.DB
	log *logger.Logger
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, log *logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:  db,
		log: log,
	}
}

// Health handles GET /health requests
func (h *HealthHandler) Health(c *gin.Context) {
	checks := make(map[string]string)

	// Check database connection
	if h.db != nil {
		if err := h.db.HealthCheck(); err != nil {
			h.log.WithError(err).Error("Database health check failed")
			checks["database"] = statusUnhealthy + ": " + err.Error()
		} else {
			checks["database"] = statusHealthy
		}
	} else {
		checks["database"] = statusNotConfigured
	}

	// Determine overall status
	status := statusHealthy
	for _, check := range checks {
		if check != statusHealthy && check != statusNotConfigured {
			status = statusUnhealthy
			break
		}
	}

	response := HealthResponse{
		Status:    status,
		Service:   "normal-form-app",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
		Checks:    checks,
	}

	// Set appropriate status code
	statusCode := http.StatusOK
	if status == statusUnhealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// LivenessProbe handles GET /health/live requests
func (h *HealthHandler) LivenessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ReadinessProbe handles GET /health/ready requests
func (h *HealthHandler) ReadinessProbe(c *gin.Context) {
	// Check if database is ready
	if h.db != nil {
		if err := h.db.HealthCheck(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "not ready",
				"reason":    "database not ready",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
