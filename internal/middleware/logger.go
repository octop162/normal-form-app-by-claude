// Package middleware provides HTTP middleware functions.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	httpStatusClientErrorStart = 400
	httpStatusClientErrorEnd   = 500
	httpStatusServerErrorStart = 500
)

// LoggerMiddleware creates a Gin middleware for logging HTTP requests
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.WithFields(map[string]interface{}{
			"timestamp":  param.TimeStamp.Format("2006-01-02T15:04:05.000Z07:00"),
			"status":     param.StatusCode,
			"latency":    param.Latency.String(),
			"client_ip":  param.ClientIP,
			"method":     param.Method,
			"path":       param.Path,
			"user_agent": param.Request.UserAgent(),
			"error":      param.ErrorMessage,
			"body_size":  param.BodySize,
		}).Info("HTTP Request")

		return ""
	})
}

// SimpleLoggerMiddleware creates a simple logger middleware
func SimpleLoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Get status code
		statusCode := c.Writer.Status()

		// Build path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log level based on status code
		logEntry := log.WithFields(map[string]interface{}{
			"status":     statusCode,
			"latency":    latency.String(),
			"client_ip":  clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
		})

		switch {
		case statusCode >= httpStatusClientErrorStart && statusCode < httpStatusClientErrorEnd:
			logEntry.Warn("Client error")
		case statusCode >= httpStatusServerErrorStart:
			logEntry.Error("Server error")
		default:
			logEntry.Info("Request completed")
		}
	}
}
