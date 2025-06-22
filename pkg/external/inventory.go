// Package external provides inventory API client functionality.
package external

import (
	"context"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	inventoryCheckEndpoint = "/api/inventory/check"
)

// InventoryClient handles inventory-related external API calls
type InventoryClient struct {
	client *Client
	log    *logger.Logger
}

// NewInventoryClient creates a new inventory API client
func NewInventoryClient(config *Config, log *logger.Logger) *InventoryClient {
	return &InventoryClient{
		client: NewClient(config, log),
		log:    log,
	}
}

// InventoryCheckRequest represents the request payload for inventory check
type InventoryCheckRequest struct {
	OptionIDs []string `json:"option_ids" validate:"required,min=1"`
}

// InventoryCheckResponse represents the response from inventory check API
type InventoryCheckResponse struct {
	Success bool              `json:"success"`
	Data    map[string]int    `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

// InventoryInfo represents inventory information for a single option
type InventoryInfo struct {
	OptionID string `json:"option_id"`
	Stock    int    `json:"stock"`
}

// CheckInventory checks the inventory levels for the specified options
func (ic *InventoryClient) CheckInventory(ctx context.Context, optionIDs []string) (map[string]int, error) {
	if len(optionIDs) == 0 {
		return nil, fmt.Errorf("option IDs cannot be empty")
	}

	// Prepare request
	req := &InventoryCheckRequest{
		OptionIDs: optionIDs,
	}

	// Make API call
	var resp InventoryCheckResponse
	err := ic.client.PostJSON(ctx, inventoryCheckEndpoint, req, &resp)
	if err != nil {
		ic.log.WithError(err).WithField("option_ids", optionIDs).Error("Failed to check inventory")
		return nil, fmt.Errorf("inventory check API call failed: %w", err)
	}

	// Validate response
	if !resp.Success {
		errMsg := "unknown error"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		ic.log.WithField("option_ids", optionIDs).WithField("api_error", errMsg).Error("Inventory API returned error")
		return nil, fmt.Errorf("inventory API error: %s", errMsg)
	}

	if resp.Data == nil {
		ic.log.WithField("option_ids", optionIDs).Error("Inventory API returned no data")
		return nil, fmt.Errorf("no inventory data received")
	}

	// Validate that all requested options are in the response
	result := make(map[string]int)
	for _, optionID := range optionIDs {
		stock, exists := resp.Data[optionID]
		if !exists {
			ic.log.WithField("option_id", optionID).Warn("Option not found in inventory response")
			// Set stock to 0 for missing options
			result[optionID] = 0
		} else {
			result[optionID] = stock
		}
	}

	ic.log.WithField("option_ids", optionIDs).WithField("inventory_result", result).Debug("Inventory check completed")
	return result, nil
}

// CheckSingleOptionInventory checks the inventory for a single option
func (ic *InventoryClient) CheckSingleOptionInventory(ctx context.Context, optionID string) (int, error) {
	if optionID == "" {
		return 0, fmt.Errorf("option ID cannot be empty")
	}

	inventory, err := ic.CheckInventory(ctx, []string{optionID})
	if err != nil {
		return 0, err
	}

	stock, exists := inventory[optionID]
	if !exists {
		return 0, fmt.Errorf("option %s not found in inventory response", optionID)
	}

	return stock, nil
}

// GetInventoryList retrieves inventory information for multiple options as a slice
func (ic *InventoryClient) GetInventoryList(ctx context.Context, optionIDs []string) ([]*InventoryInfo, error) {
	inventory, err := ic.CheckInventory(ctx, optionIDs)
	if err != nil {
		return nil, err
	}

	result := make([]*InventoryInfo, 0, len(inventory))
	for optionID, stock := range inventory {
		result = append(result, &InventoryInfo{
			OptionID: optionID,
			Stock:    stock,
		})
	}

	return result, nil
}