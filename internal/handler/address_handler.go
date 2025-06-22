// Package handler provides HTTP handlers for address management.
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/service"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// AddressHandler handles address-related HTTP requests
type AddressHandler struct {
	addressService service.AddressService
	log            *logger.Logger
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressService service.AddressService, log *logger.Logger) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		log:            log,
	}
}

// SearchAddress handles GET /api/v1/address/search
func (h *AddressHandler) SearchAddress(c *gin.Context) {
	var req dto.AddressSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind address search request")
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

	// Search address by postal code
	resp, err := h.addressService.SearchByPostalCode(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to search address")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeAddressSearchFailed,
				Message: "Failed to search address",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// CheckRegion handles POST /api/v1/region/check
func (h *AddressHandler) CheckRegion(c *gin.Context) {
	var req dto.RegionCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Error("Failed to bind region check request")
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

	// Check region restrictions
	resp, err := h.addressService.CheckRegionRestrictions(c.Request.Context(), &req)
	if err != nil {
		h.log.WithError(err).Error("Failed to check region restrictions")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeRegionCheckFailed,
				Message: "Failed to check region restrictions",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetPrefectures handles GET /api/v1/prefectures
func (h *AddressHandler) GetPrefectures(c *gin.Context) {
	// Get prefectures
	resp, err := h.addressService.GetPrefectures(c.Request.Context())
	if err != nil {
		h.log.WithError(err).Error("Failed to get prefectures")
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeInternalError,
				Message: "Failed to retrieve prefectures",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    resp,
	})
}

// GetPrefecture handles GET /api/v1/prefectures/:name
func (h *AddressHandler) GetPrefecture(c *gin.Context) {
	prefectureName := c.Param("name")
	if prefectureName == "" {
		h.log.Error("Missing prefecture name")
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Error: &dto.APIError{
				Code:    ErrorCodeMissingPrefectureName,
				Message: "Prefecture name is required",
			},
		})
		return
	}

	// Get prefecture by name
	resp, err := h.addressService.GetPrefectureByName(c.Request.Context(), prefectureName)
	if err != nil {
		h.log.WithError(err).WithField("prefecture_name", prefectureName).Error("Failed to get prefecture")

		statusCode := http.StatusInternalServerError
		errorCode := ErrorCodeInternalError

		if isNotFoundError(err) {
			statusCode = http.StatusNotFound
			errorCode = ErrorCodePrefectureNotFound
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
