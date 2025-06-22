// Package handler provides HTTP handlers for user management.
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
	log         *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind user create request")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidRequest,
				Message: "Invalid request format",
				Details: map[string]string{"bind_error": err.Error()},
			},
		})
		return
	}

	// Create user
	resp, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to create user")

		// Check for specific error types
		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		switch {
		case isValidationError(err):
			statusCode = http.StatusBadRequest
			errorCode = ErrorCodeValidationError
		case isDuplicateError(err):
			statusCode = http.StatusConflict
			errorCode = ErrorCodeDuplicateError
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    errorCode,
				Message: err.Error(),
			},
		})
		return
	}

	h.log.WithField("user_id", resp.ID).Info("User created successfully")
	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// ValidateUser handles POST /api/v1/users/validate
func (h *UserHandler) ValidateUser(c *gin.Context) {
	var req dto.UserValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind user validate request")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidRequest,
				Message: "Invalid request format",
				Details: map[string]string{"bind_error": err.Error()},
			},
		})
		return
	}

	// Validate user data
	resp, err := h.userService.ValidateUserData(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to validate user data")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeValidationError,
				Message: "Validation process failed",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetUser handles GET /api/v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.WithError(err).WithField("id_param", idParam).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidUserID,
				Message: "User ID must be a valid integer",
			},
		})
		return
	}

	// Get user
	resp, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to get user")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeUserNotFound
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    errorCode,
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.WithError(err).WithField("id_param", idParam).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidUserID,
				Message: "User ID must be a valid integer",
			},
		})
		return
	}

	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind user update request")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidRequest,
				Message: "Invalid request format",
				Details: map[string]string{"bind_error": err.Error()},
			},
		})
		return
	}

	// Update user
	resp, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to update user")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isValidationError(err) {
			statusCode = http.StatusBadRequest
			errorCode = ErrorCodeValidationError
		} else if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeUserNotFound
		} else if isDuplicateError(err) {
			statusCode = http.StatusConflict
			errorCode = ErrorCodeDuplicateError
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    errorCode,
				Message: err.Error(),
			},
		})
		return
	}

	h.log.WithField("user_id", userID).Info("User updated successfully")
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.WithError(err).WithField("id_param", idParam).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidUserID,
				Message: "User ID must be a valid integer",
			},
		})
		return
	}

	// Delete user
	err = h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.log.WithError(err).WithField("user_id", userID).Error("Failed to delete user")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeUserNotFound
		}

		c.JSON(statusCode, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    errorCode,
				Message: err.Error(),
			},
		})
		return
	}

	h.log.WithField("user_id", userID).Info("User deleted successfully")
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "User deleted successfully"},
	})
}
