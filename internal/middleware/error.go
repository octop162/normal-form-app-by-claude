package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ErrorHandlerMiddleware handles errors and panics
func ErrorHandlerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.WithFields(map[string]interface{}{
				"error":      err,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Error("Panic recovered")

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal Server Error",
				Message: "An unexpected error occurred",
				Code:    http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		if err, ok := recovered.(error); ok {
			log.WithError(err).WithFields(map[string]interface{}{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Error("Panic recovered")

			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal Server Error",
				Message: "An unexpected error occurred",
				Code:    http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		// Default case
		log.WithFields(map[string]interface{}{
			"recovered":  recovered,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Error("Unknown panic recovered")

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: "An unexpected error occurred",
			Code:    http.StatusInternalServerError,
		})
		c.Abort()
	})
}

// NotFoundMiddleware handles 404 errors
func NotFoundMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Not Found",
			Message: "The requested resource was not found",
			Code:    http.StatusNotFound,
		})
	}
}

// MethodNotAllowedMiddleware handles 405 errors
func MethodNotAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "Method Not Allowed",
			Message: "The requested method is not allowed for this resource",
			Code:    http.StatusMethodNotAllowed,
		})
	}
}
