// Package handler provides HTTP handlers for session management.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// SessionHandler handles session-related HTTP requests
type SessionHandler struct {
	sessionService service.SessionService
	log            *logger.Logger
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService service.SessionService, log *logger.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		log:            log,
	}
}

// CreateSession handles POST /api/v1/sessions
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req dto.SessionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind session create request")
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

	// Create session
	resp, err := h.sessionService.CreateSession(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to create session")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeSessionCreateFailed,
				Message: "Failed to create session",
			},
		})
		return
	}

	h.log.WithField("session_id", resp.SessionID).Info("Session created successfully")
	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetSession handles GET /api/v1/sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		h.log.Error("Missing session ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingSessionID,
				Message: "Session ID is required",
			},
		})
		return
	}

	// Get session
	resp, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		h.log.WithError(err).WithField("session_id", sessionID).Error("Failed to get session")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) || isExpiredError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeSessionNotFound
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

// UpdateSession handles PUT /api/v1/sessions/:id
func (h *SessionHandler) UpdateSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		h.log.Error("Missing session ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingSessionID,
				Message: "Session ID is required",
			},
		})
		return
	}

	var req dto.SessionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind session update request")
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

	// Update session
	resp, err := h.sessionService.UpdateSession(c.Request.Context(), sessionID, &req)
	if err != nil {
		h.log.WithError(err).WithField("session_id", sessionID).Error("Failed to update session")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) || isExpiredError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeSessionNotFound
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

	h.log.WithField("session_id", sessionID).Info("Session updated successfully")
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// DeleteSession handles DELETE /api/v1/sessions/:id
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		h.log.Error("Missing session ID")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingSessionID,
				Message: "Session ID is required",
			},
		})
		return
	}

	// Delete session
	resp, err := h.sessionService.DeleteSession(c.Request.Context(), sessionID)
	if err != nil {
		h.log.WithError(err).WithField("session_id", sessionID).Error("Failed to delete session")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeSessionNotFound
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

	h.log.WithField("session_id", sessionID).Info("Session deleted successfully")
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}
