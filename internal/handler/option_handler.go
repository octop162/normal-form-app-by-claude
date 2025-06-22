// Package handler provides HTTP handlers for option management.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// OptionHandler handles option-related HTTP requests
type OptionHandler struct {
	optionService service.OptionService
	log           *logger.Logger
}

// NewOptionHandler creates a new option handler
func NewOptionHandler(optionService service.OptionService, log *logger.Logger) *OptionHandler {
	return &OptionHandler{
		optionService: optionService,
		log:           log,
	}
}

// GetOptions handles GET /api/v1/options
func (h *OptionHandler) GetOptions(c *gin.Context) {
	var req dto.OptionsGetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind options get request")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInvalidRequest,
				Message: "Invalid query parameters",
				Details: map[string]string{"bind_error": err.Error()},
			},
		})
		return
	}

	// Get available options
	resp, err := h.optionService.GetAvailableOptions(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to get available options")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInternalError,
				Message: "Failed to retrieve options",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// CheckInventory handles POST /api/v1/options/check-inventory
func (h *OptionHandler) CheckInventory(c *gin.Context) {
	var req dto.InventoryCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind inventory check request")
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

	// Check inventory
	resp, err := h.optionService.CheckInventory(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to check inventory")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInventoryCheckFailed,
				Message: "Failed to check inventory levels",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetOption handles GET /api/v1/options/:type
func (h *OptionHandler) GetOption(c *gin.Context) {
	optionType := c.Param("type")
	if optionType == "" {
		h.log.Error("Missing option type")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingOptionType,
				Message: "Option type is required",
			},
		})
		return
	}

	// Get option by type
	resp, err := h.optionService.GetOptionByType(c.Request.Context(), optionType)
	if err != nil {
		h.log.WithError(err).WithField("option_type", optionType).Error("Failed to get option")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodeOptionNotFound
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
