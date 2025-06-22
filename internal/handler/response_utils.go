// Package handler provides response utilities for HTTP handlers.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// respondWithError sends an error response
func respondWithError(c *gin.Context, statusCode int, errorCode, message string, log *logger.Logger, err error) {
	if log != nil && err != nil {
		log.WithError(err).Error(message)
	}

	c.JSON(statusCode, dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    errorCode,
			Message: message,
		},
	})
}

// respondWithBindError sends a bind error response
func respondWithBindError(c *gin.Context, err error, log *logger.Logger, operation string) {
	if log != nil {
		log.WithError(err).Errorf("Failed to bind %s request", operation)
	}

	c.JSON(http.StatusBadRequest, dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    ErrorCodeInvalidRequest,
			Message: MessageInvalidRequest,
			Details: map[string]string{"bind_error": err.Error()},
		},
	})
}

// respondWithSuccess sends a success response
func respondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, dto.APIResponse{
		Success: true,
		Data:    data,
	})
}

// handleServiceError determines the appropriate error response based on error type
func handleServiceError(c *gin.Context, err error, log *logger.Logger, operation string, notFoundCode string) {
	statusCode := http.StatusInternalServerError
	errorCode := ErrorCodeInternalError

	switch {
	case isValidationError(err):
		statusCode = http.StatusBadRequest
		errorCode = ErrorCodeValidationError
	case isNotFoundError(err):
		statusCode = http.StatusNotFound
		errorCode = notFoundCode
	case isDuplicateError(err):
		statusCode = http.StatusConflict
		errorCode = ErrorCodeDuplicateError
	case isExpiredError(err):
		statusCode = http.StatusNotFound
		errorCode = notFoundCode
	}

	if log != nil {
		log.WithError(err).Errorf("Failed to %s", operation)
	}

	c.JSON(statusCode, dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    errorCode,
			Message: err.Error(),
		},
	})
}

// validatePathParam validates that a path parameter is not empty
func validatePathParam(c *gin.Context, paramName, paramValue, errorCode, errorMessage string, log *logger.Logger) bool {
	if paramValue == "" {
		if log != nil {
			log.Errorf("Missing %s", paramName)
		}
		respondWithError(c, http.StatusBadRequest, errorCode, errorMessage, nil, nil)
		return false
	}
	return true
}
