// Package handler provides HTTP handlers for plan management.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// PlanHandler handles plan-related HTTP requests
type PlanHandler struct {
	planService service.PlanService
	log         *logger.Logger
}

// NewPlanHandler creates a new plan handler
func NewPlanHandler(planService service.PlanService, log *logger.Logger) *PlanHandler {
	return &PlanHandler{
		planService: planService,
		log:         log,
	}
}

// GetPlans handles GET /api/v1/plans
func (h *PlanHandler) GetPlans(c *gin.Context) {
	// Get available plans
	resp, err := h.planService.GetAvailablePlans(c.Request.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to get available plans")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInternalError,
				Message: "Failed to retrieve plans",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetPlan handles GET /api/v1/plans/:type
func (h *PlanHandler) GetPlan(c *gin.Context) {
	planType := c.Param("type")
	if planType == "" {
		h.log.Error("Missing plan type")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingPlanType,
				Message: "Plan type is required",
			},
		})
		return
	}

	// Get plan by type
	resp, err := h.planService.GetPlanByType(c.Request.Context(), planType)
	if err != nil {
		h.log.WithError(err).WithField("plan_type", planType).Error("Failed to get plan")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodePlanNotFound
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
