// Package service provides option management business logic.
package service

import (
	"context"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/internal/repository"
	"github.com/octop162/normal-form-app-by-claude/pkg/external"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	// Mock inventory levels for testing
	mockInventoryAA       = 10
	mockInventoryAB       = 25
	defaultInventoryLevel = 5
)

// OptionService defines the interface for option business logic
type OptionService interface {
	GetAvailableOptions(ctx context.Context, req *dto.OptionsGetRequest) (*dto.OptionsGetResponse, error)
	CheckInventory(ctx context.Context, req *dto.InventoryCheckRequest) (*dto.InventoryCheckResponse, error)
	GetOptionByType(ctx context.Context, optionType string) (*dto.OptionResponse, error)
	GetAllOptions(ctx context.Context) (*dto.OptionsGetResponse, error)
}

// optionService implements OptionService
type optionService struct {
	optionRepo  repository.OptionRepository
	externalAPI *external.Manager
	log         *logger.Logger
}

// NewOptionService creates a new option service
func NewOptionService(
	optionRepo repository.OptionRepository,
	externalAPI *external.Manager,
	log *logger.Logger,
) OptionService {
	return &optionService{
		optionRepo:  optionRepo,
		externalAPI: externalAPI,
		log:         log,
	}
}

// GetAvailableOptions retrieves options available for a specific plan type
func (s *optionService) GetAvailableOptions(
	ctx context.Context, req *dto.OptionsGetRequest,
) (*dto.OptionsGetResponse, error) {
	var options []*model.OptionMaster
	var err error

	if req.PlanType != "" {
		// Get options compatible with the specified plan type
		options, err = s.optionRepo.GetByPlanType(ctx, req.PlanType)
		if err != nil {
			s.log.WithError(err).WithField("plan_type", req.PlanType).Error("Failed to get options by plan type")
			return nil, fmt.Errorf("failed to get options by plan type: %w", err)
		}

		// TODO: Apply region restrictions if region is specified
		if req.Region != "" {
			options = s.filterOptionsByRegion(options, req.Region)
		}
	} else {
		// Get all active options
		options, err = s.optionRepo.GetActiveOptions(ctx)
		if err != nil {
			s.log.WithError(err).Error("Failed to get all active options")
			return nil, fmt.Errorf("failed to get all active options: %w", err)
		}
	}

	// Convert to response DTOs
	optionResponses := make([]dto.OptionResponse, len(options))
	for i, option := range options {
		optionResponses[i] = s.convertOptionToResponse(option)
	}

	return &dto.OptionsGetResponse{
		Options: optionResponses,
	}, nil
}

// CheckInventory checks inventory levels for specified option types
func (s *optionService) CheckInventory(
	ctx context.Context, req *dto.InventoryCheckRequest,
) (*dto.InventoryCheckResponse, error) {
	inventory := make(map[string]int)

	// Try external inventory API first if available
	if s.externalAPI != nil && s.externalAPI.InventoryClient() != nil {
		externalInventory, err := s.externalAPI.InventoryClient().CheckInventory(ctx, req.OptionTypes)
		if err != nil {
			s.log.WithError(err).WithField("option_types", req.OptionTypes).Warn("External inventory API failed, falling back to local logic")
		} else {
			// Validate options exist in local database and are active
			for optionType, stock := range externalInventory {
				option, err := s.optionRepo.GetByOptionType(ctx, optionType)
				if err != nil || !option.IsActive {
					inventory[optionType] = 0
				} else {
					inventory[optionType] = stock
				}
			}
			return &dto.InventoryCheckResponse{
				Inventory: inventory,
			}, nil
		}
	}

	// Fallback to local logic
	for _, optionType := range req.OptionTypes {
		option, err := s.optionRepo.GetByOptionType(ctx, optionType)
		if err != nil {
			s.log.WithError(err).WithField("option_type", optionType).Error("Failed to get option")
			// Set inventory to 0 for non-existent options
			inventory[optionType] = 0
			continue
		}

		if !option.IsActive {
			// Inactive options have 0 inventory
			inventory[optionType] = 0
			continue
		}

		// Use mock data as fallback
		inventoryLevel := s.getMockInventoryLevel(optionType)
		inventory[optionType] = inventoryLevel
	}

	return &dto.InventoryCheckResponse{
		Inventory: inventory,
	}, nil
}

// GetOptionByType retrieves a specific option by its type
func (s *optionService) GetOptionByType(ctx context.Context, optionType string) (*dto.OptionResponse, error) {
	option, err := s.optionRepo.GetByOptionType(ctx, optionType)
	if err != nil {
		s.log.WithError(err).WithField("option_type", optionType).Error("Failed to get option by type")
		return nil, fmt.Errorf("failed to get option by type: %w", err)
	}

	response := s.convertOptionToResponse(option)
	return &response, nil
}

// GetAllOptions retrieves all options
func (s *optionService) GetAllOptions(ctx context.Context) (*dto.OptionsGetResponse, error) {
	options, err := s.optionRepo.GetAll(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to get all options")
		return nil, fmt.Errorf("failed to get all options: %w", err)
	}

	// Convert to response DTOs
	optionResponses := make([]dto.OptionResponse, len(options))
	for i, option := range options {
		optionResponses[i] = s.convertOptionToResponse(option)
	}

	return &dto.OptionsGetResponse{
		Options: optionResponses,
	}, nil
}

// convertOptionToResponse converts option model to response DTO
func (s *optionService) convertOptionToResponse(option *model.OptionMaster) dto.OptionResponse {
	description := ""
	if option.Description != nil {
		description = *option.Description
	}

	return dto.OptionResponse{
		ID:                option.ID,
		OptionType:        option.OptionType,
		OptionName:        option.OptionName,
		Description:       description,
		PlanCompatibility: option.PlanCompatibility,
		IsActive:          option.IsActive,
	}
}

// filterOptionsByRegion filters options based on region restrictions
// TODO: Implement actual region-based filtering logic
func (s *optionService) filterOptionsByRegion(options []*model.OptionMaster, region string) []*model.OptionMaster {
	// For now, return all options without filtering
	// In production, this would call external region restriction API
	s.log.WithField("region", region).Debug("Region-based filtering not yet implemented")
	return options
}

// getMockInventoryLevel returns mock inventory levels for testing
// TODO: Replace with actual external API call
func (s *optionService) getMockInventoryLevel(optionType string) int {
	// Mock inventory data for testing
	mockInventory := map[string]int{
		"AA": mockInventoryAA,
		"BB": 0, // Out of stock
		"AB": mockInventoryAB,
	}

	if level, exists := mockInventory[optionType]; exists {
		return level
	}

	// Default inventory for unknown options
	return defaultInventoryLevel
}
