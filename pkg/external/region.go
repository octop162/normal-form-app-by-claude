// Package external provides region restriction API client functionality.
package external

import (
	"context"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	regionCheckEndpoint = "/api/region/check"
)

// RegionClient handles region restriction-related external API calls
type RegionClient struct {
	client *Client
	log    *logger.Logger
}

// NewRegionClient creates a new region API client
func NewRegionClient(config *Config, log *logger.Logger) *RegionClient {
	return &RegionClient{
		client: NewClient(config, log),
		log:    log,
	}
}

// RegionCheckRequest represents the request payload for region restriction check
type RegionCheckRequest struct {
	Prefecture string   `json:"prefecture" validate:"required"`
	City       string   `json:"city" validate:"required"`
	OptionIDs  []string `json:"option_ids" validate:"required,min=1"`
}

// RegionCheckResponse represents the response from region check API
type RegionCheckResponse struct {
	Success bool           `json:"success"`
	Data    map[string]bool `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// RegionRestrictionInfo represents region restriction information for a single option
type RegionRestrictionInfo struct {
	OptionID   string `json:"option_id"`
	IsAllowed  bool   `json:"is_allowed"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
}

// CheckRegionRestrictions checks if the specified options are allowed in the given region
func (rc *RegionClient) CheckRegionRestrictions(ctx context.Context, prefecture, city string, optionIDs []string) (map[string]bool, error) {
	if prefecture == "" {
		return nil, fmt.Errorf("prefecture cannot be empty")
	}
	if city == "" {
		return nil, fmt.Errorf("city cannot be empty")
	}
	if len(optionIDs) == 0 {
		return nil, fmt.Errorf("option IDs cannot be empty")
	}

	// Prepare request
	req := &RegionCheckRequest{
		Prefecture: prefecture,
		City:       city,
		OptionIDs:  optionIDs,
	}

	// Make API call
	var resp RegionCheckResponse
	err := rc.client.PostJSON(ctx, regionCheckEndpoint, req, &resp)
	if err != nil {
		rc.log.WithError(err).
			WithField("prefecture", prefecture).
			WithField("city", city).
			WithField("option_ids", optionIDs).
			Error("Failed to check region restrictions")
		return nil, fmt.Errorf("region check API call failed: %w", err)
	}

	// Validate response
	if !resp.Success {
		errMsg := "unknown error"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		rc.log.WithField("prefecture", prefecture).
			WithField("city", city).
			WithField("option_ids", optionIDs).
			WithField("api_error", errMsg).
			Error("Region API returned error")
		return nil, fmt.Errorf("region API error: %s", errMsg)
	}

	if resp.Data == nil {
		rc.log.WithField("prefecture", prefecture).
			WithField("city", city).
			WithField("option_ids", optionIDs).
			Error("Region API returned no data")
		return nil, fmt.Errorf("no region restriction data received")
	}

	// Validate that all requested options are in the response
	result := make(map[string]bool)
	for _, optionID := range optionIDs {
		isAllowed, exists := resp.Data[optionID]
		if !exists {
			rc.log.WithField("option_id", optionID).
				WithField("prefecture", prefecture).
				WithField("city", city).
				Warn("Option not found in region restriction response")
			// Default to not allowed for missing options (safe default)
			result[optionID] = false
		} else {
			result[optionID] = isAllowed
		}
	}

	rc.log.WithField("prefecture", prefecture).
		WithField("city", city).
		WithField("option_ids", optionIDs).
		WithField("restriction_result", result).
		Debug("Region restriction check completed")
	return result, nil
}

// CheckSingleOptionRegionRestriction checks if a single option is allowed in the given region
func (rc *RegionClient) CheckSingleOptionRegionRestriction(ctx context.Context, prefecture, city, optionID string) (bool, error) {
	if optionID == "" {
		return false, fmt.Errorf("option ID cannot be empty")
	}

	restrictions, err := rc.CheckRegionRestrictions(ctx, prefecture, city, []string{optionID})
	if err != nil {
		return false, err
	}

	isAllowed, exists := restrictions[optionID]
	if !exists {
		return false, fmt.Errorf("option %s not found in region restriction response", optionID)
	}

	return isAllowed, nil
}

// GetRegionRestrictionList retrieves region restriction information for multiple options as a slice
func (rc *RegionClient) GetRegionRestrictionList(ctx context.Context, prefecture, city string, optionIDs []string) ([]*RegionRestrictionInfo, error) {
	restrictions, err := rc.CheckRegionRestrictions(ctx, prefecture, city, optionIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*RegionRestrictionInfo, 0, len(restrictions))
	for optionID, isAllowed := range restrictions {
		result = append(result, &RegionRestrictionInfo{
			OptionID:   optionID,
			IsAllowed:  isAllowed,
			Prefecture: prefecture,
			City:       city,
		})
	}

	return result, nil
}

// GetAllowedOptions filters and returns only the options that are allowed in the given region
func (rc *RegionClient) GetAllowedOptions(ctx context.Context, prefecture, city string, optionIDs []string) ([]string, error) {
	restrictions, err := rc.CheckRegionRestrictions(ctx, prefecture, city, optionIDs)
	if err != nil {
		return nil, err
	}

	var allowedOptions []string
	for optionID, isAllowed := range restrictions {
		if isAllowed {
			allowedOptions = append(allowedOptions, optionID)
		}
	}

	return allowedOptions, nil
}

// GetRestrictedOptions filters and returns only the options that are restricted in the given region
func (rc *RegionClient) GetRestrictedOptions(ctx context.Context, prefecture, city string, optionIDs []string) ([]string, error) {
	restrictions, err := rc.CheckRegionRestrictions(ctx, prefecture, city, optionIDs)
	if err != nil {
		return nil, err
	}

	var restrictedOptions []string
	for optionID, isAllowed := range restrictions {
		if !isAllowed {
			restrictedOptions = append(restrictedOptions, optionID)
		}
	}

	return restrictedOptions, nil
}