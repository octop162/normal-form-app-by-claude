package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents error codes for the application
type ErrorCode string

const (
	// Generic error codes
	ErrorCodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrorCodeNotFoundGeneric ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrorCodeConflict        ErrorCode = "CONFLICT"
	ErrorCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"

	// Validation error codes
	ErrorCodeValidationFailed      ErrorCode = "VALIDATION_FAILED"
	ErrorCodeRequiredFieldMissing  ErrorCode = "REQUIRED_FIELD_MISSING"
	ErrorCodeInvalidFormat         ErrorCode = "INVALID_FORMAT"
	ErrorCodeValueTooLong          ErrorCode = "VALUE_TOO_LONG"
	ErrorCodeValueTooShort         ErrorCode = "VALUE_TOO_SHORT"
	ErrorCodeInvalidEmail          ErrorCode = "INVALID_EMAIL"
	ErrorCodeInvalidPhoneNumber    ErrorCode = "INVALID_PHONE_NUMBER"
	ErrorCodeInvalidPostalCode     ErrorCode = "INVALID_POSTAL_CODE"
	ErrorCodeEmailConfirmationFail ErrorCode = "EMAIL_CONFIRMATION_FAILED"

	// Business logic error codes
	ErrorCodeUserAlreadyExists     ErrorCode = "USER_ALREADY_EXISTS"
	ErrorCodeUserNotFound          ErrorCode = "USER_NOT_FOUND"
	ErrorCodeSessionExpired        ErrorCode = "SESSION_EXPIRED"
	ErrorCodeSessionNotFoundError  ErrorCode = "SESSION_NOT_FOUND"
	ErrorCodeInvalidSessionData    ErrorCode = "INVALID_SESSION_DATA"
	ErrorCodeInventoryNotAvailable ErrorCode = "INVENTORY_NOT_AVAILABLE"
	ErrorCodeRegionNotSupported    ErrorCode = "REGION_NOT_SUPPORTED"
	ErrorCodeOptionNotAvailable    ErrorCode = "OPTION_NOT_AVAILABLE"
	ErrorCodePlanNotFoundError     ErrorCode = "PLAN_NOT_FOUND"
	ErrorCodeAddressNotFound       ErrorCode = "ADDRESS_NOT_FOUND"

	// External API error codes
	ErrorCodeExternalAPIError     ErrorCode = "EXTERNAL_API_ERROR"
	ErrorCodeInventoryAPIError    ErrorCode = "INVENTORY_API_ERROR"
	ErrorCodeAddressAPIError      ErrorCode = "ADDRESS_API_ERROR"
	ErrorCodeRegionAPIError       ErrorCode = "REGION_API_ERROR"
	ErrorCodeExternalAPITimeout   ErrorCode = "EXTERNAL_API_TIMEOUT"
	ErrorCodeExternalAPIRateLimit ErrorCode = "EXTERNAL_API_RATE_LIMIT"

	// Security error codes
	ErrorCodeCSRFTokenMissing     ErrorCode = "CSRF_TOKEN_MISSING"
	ErrorCodeCSRFTokenInvalid     ErrorCode = "CSRF_TOKEN_INVALID"
	ErrorCodeRateLimitExceeded    ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeSuspiciousActivity   ErrorCode = "SUSPICIOUS_ACTIVITY"
	ErrorCodeUnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"

	// Database error codes
	ErrorCodeDatabaseError       ErrorCode = "DATABASE_ERROR"
	ErrorCodeDatabaseTimeout     ErrorCode = "DATABASE_TIMEOUT"
	ErrorCodeDatabaseConnection  ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrorCodeDuplicateEntry      ErrorCode = "DUPLICATE_ENTRY"
	ErrorCodeConstraintViolation ErrorCode = "CONSTRAINT_VIOLATION"
)

// AppError represents application-specific errors
type AppError struct {
	Code       ErrorCode         `json:"code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"-"`
	Err        error             `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// NewValidationError creates a validation error with field details
func NewValidationError(field string, message string) *AppError {
	return &AppError{
		Code:       ErrorCodeValidationFailed,
		Message:    "入力内容に不備があります",
		StatusCode: http.StatusBadRequest,
		Details: map[string]string{
			field: message,
		},
	}
}

// NewBusinessLogicError creates a business logic error
func NewBusinessLogicError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewExternalAPIError creates an external API error
func NewExternalAPIError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
		Err:        err,
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   *ErrorDetail      `json:"error"`
	Meta    *ErrorMeta        `json:"meta,omitempty"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// ErrorMeta represents metadata for error tracking
type ErrorMeta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Path      string `json:"path,omitempty"`
	Method    string `json:"method,omitempty"`
}

// HandleError handles application errors and returns appropriate HTTP responses
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	var statusCode int
	var errorDetail *ErrorDetail

	// Check if it's an AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
		statusCode = e.StatusCode
		errorDetail = &ErrorDetail{
			Code:    e.Code,
			Message: e.Message,
			Details: e.Details,
		}
	} else {
		// Generic error handling
		statusCode = http.StatusInternalServerError
		errorDetail = &ErrorDetail{
			Code:    ErrorCodeInternalServer,
			Message: "内部サーバーエラーが発生しました",
		}
	}

	// Create error response
	response := ErrorResponse{
		Success: false,
		Error:   errorDetail,
		Meta: &ErrorMeta{
			RequestID: c.GetHeader("X-Request-ID"),
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
		},
	}

	// Log error for debugging
	if appErr != nil && appErr.Err != nil {
		c.Set("error", appErr.Err)
	}

	c.JSON(statusCode, response)
}

// HandleValidationErrors handles multiple validation errors
func HandleValidationErrors(c *gin.Context, errors map[string]string) {
	response := ErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:    ErrorCodeValidationFailed,
			Message: "入力内容に不備があります",
			Details: errors,
		},
		Meta: &ErrorMeta{
			RequestID: c.GetHeader("X-Request-ID"),
			Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
			Path:      c.Request.URL.Path,
			Method:    c.Request.Method,
		},
	}

	c.JSON(http.StatusBadRequest, response)
}

// HandleSuccessResponse handles successful responses
func HandleSuccessResponse(c *gin.Context, data interface{}) {
	response := gin.H{
		"success": true,
		"data":    data,
	}

	c.JSON(http.StatusOK, response)
}

// HandleCreatedResponse handles successful creation responses
func HandleCreatedResponse(c *gin.Context, data interface{}) {
	response := gin.H{
		"success": true,
		"data":    data,
	}

	c.JSON(http.StatusCreated, response)
}

// HandleNoContentResponse handles successful no-content responses
func HandleNoContentResponse(c *gin.Context) {
	response := gin.H{
		"success": true,
	}

	c.JSON(http.StatusOK, response)
}

// PanicRecovery handles panics and converts them to errors
func PanicRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var message string
		switch v := recovered.(type) {
		case string:
			message = v
		case error:
			message = v.Error()
		default:
			message = "予期しないエラーが発生しました"
		}

		appErr := NewAppError(
			ErrorCodeInternalServer,
			message,
			http.StatusInternalServerError,
			nil,
		)

		HandleError(c, appErr)
		c.Abort()
	})
}