// Package external provides a unified manager for all external API clients.
package external

import (
	"context"

	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// Manager provides a unified interface for all external API clients
type Manager struct {
	inventory *InventoryClient
	region    *RegionClient
	address   *AddressClient
	log       *logger.Logger
}

// ManagerConfig holds configuration for all external API clients
type ManagerConfig struct {
	InventoryAPI *Config `json:"inventory_api"`
	RegionAPI    *Config `json:"region_api"`
	AddressAPI   *Config `json:"address_api"`
}

// NewManager creates a new external API manager with all clients
func NewManager(config *ManagerConfig, log *logger.Logger) *Manager {
	var inventory *InventoryClient
	var region *RegionClient
	var address *AddressClient

	if config.InventoryAPI != nil {
		inventory = NewInventoryClient(config.InventoryAPI, log)
	}

	if config.RegionAPI != nil {
		region = NewRegionClient(config.RegionAPI, log)
	}

	if config.AddressAPI != nil {
		address = NewAddressClient(config.AddressAPI, log)
	}

	return &Manager{
		inventory: inventory,
		region:    region,
		address:   address,
		log:       log,
	}
}

// InventoryClient returns the inventory API client
func (m *Manager) InventoryClient() *InventoryClient {
	return m.inventory
}

// RegionClient returns the region API client
func (m *Manager) RegionClient() *RegionClient {
	return m.region
}

// AddressClient returns the address API client
func (m *Manager) AddressClient() *AddressClient {
	return m.address
}

// CheckOptionAvailability checks both inventory and region restrictions for options
func (m *Manager) CheckOptionAvailability(ctx context.Context, prefecture, city string, optionIDs []string) (*OptionAvailabilityResult, error) {
	result := &OptionAvailabilityResult{
		OptionResults: make(map[string]*OptionAvailability),
	}

	// Check inventory if client is available
	var inventoryMap map[string]int
	if m.inventory != nil {
		var err error
		inventoryMap, err = m.inventory.CheckInventory(ctx, optionIDs)
		if err != nil {
			m.log.WithError(err).WithField("option_ids", optionIDs).Warn("Failed to check inventory, continuing without inventory data")
			// Continue without inventory data - don't fail the entire operation
		}
	}

	// Check region restrictions if client is available
	var regionMap map[string]bool
	if m.region != nil && prefecture != "" && city != "" {
		var err error
		regionMap, err = m.region.CheckRegionRestrictions(ctx, prefecture, city, optionIDs)
		if err != nil {
			m.log.WithError(err).
				WithField("prefecture", prefecture).
				WithField("city", city).
				WithField("option_ids", optionIDs).
				Warn("Failed to check region restrictions, continuing without region data")
			// Continue without region data - don't fail the entire operation
		}
	}

	// Combine results
	for _, optionID := range optionIDs {
		availability := &OptionAvailability{
			OptionID: optionID,
		}

		// Set inventory data
		if inventoryMap != nil {
			if stock, exists := inventoryMap[optionID]; exists {
				availability.Stock = &stock
				availability.HasStock = stock > 0
			}
		}

		// Set region restriction data
		if regionMap != nil {
			if isAllowed, exists := regionMap[optionID]; exists {
				availability.IsRegionAllowed = &isAllowed
			}
		}

		// Determine overall availability
		availability.IsAvailable = availability.HasStock && (availability.IsRegionAllowed == nil || *availability.IsRegionAllowed)

		result.OptionResults[optionID] = availability
	}

	return result, nil
}

// OptionAvailabilityResult represents the combined availability check result
type OptionAvailabilityResult struct {
	OptionResults map[string]*OptionAvailability `json:"option_results"`
}

// OptionAvailability represents the availability status of a single option
type OptionAvailability struct {
	OptionID        string `json:"option_id"`
	Stock           *int   `json:"stock,omitempty"`
	HasStock        bool   `json:"has_stock"`
	IsRegionAllowed *bool  `json:"is_region_allowed,omitempty"`
	IsAvailable     bool   `json:"is_available"`
}

// GetAvailableOptions returns only the options that are available
func (r *OptionAvailabilityResult) GetAvailableOptions() []string {
	var available []string
	for optionID, availability := range r.OptionResults {
		if availability.IsAvailable {
			available = append(available, optionID)
		}
	}
	return available
}

// GetUnavailableOptions returns only the options that are not available
func (r *OptionAvailabilityResult) GetUnavailableOptions() []string {
	var unavailable []string
	for optionID, availability := range r.OptionResults {
		if !availability.IsAvailable {
			unavailable = append(unavailable, optionID)
		}
	}
	return unavailable
}

// GetOutOfStockOptions returns options that are out of stock
func (r *OptionAvailabilityResult) GetOutOfStockOptions() []string {
	var outOfStock []string
	for optionID, availability := range r.OptionResults {
		if !availability.HasStock {
			outOfStock = append(outOfStock, optionID)
		}
	}
	return outOfStock
}

// GetRegionRestrictedOptions returns options that are restricted in the region
func (r *OptionAvailabilityResult) GetRegionRestrictedOptions() []string {
	var restricted []string
	for optionID, availability := range r.OptionResults {
		if availability.IsRegionAllowed != nil && !*availability.IsRegionAllowed {
			restricted = append(restricted, optionID)
		}
	}
	return restricted
}

// HealthCheck performs health checks on all configured external APIs
func (m *Manager) HealthCheck(ctx context.Context) *HealthCheckResult {
	result := &HealthCheckResult{
		Services: make(map[string]*ServiceHealth),
	}

	// Check inventory API
	if m.inventory != nil {
		health := &ServiceHealth{Name: "inventory"}
		_, err := m.inventory.CheckInventory(ctx, []string{"TEST"})
		if err != nil {
			health.Status = "unhealthy"
			health.Error = err.Error()
		} else {
			health.Status = "healthy"
		}
		result.Services["inventory"] = health
	}

	// Check region API
	if m.region != nil {
		health := &ServiceHealth{Name: "region"}
		_, err := m.region.CheckRegionRestrictions(ctx, "東京都", "渋谷区", []string{"TEST"})
		if err != nil {
			health.Status = "unhealthy"
			health.Error = err.Error()
		} else {
			health.Status = "healthy"
		}
		result.Services["region"] = health
	}

	// Check address API
	if m.address != nil {
		health := &ServiceHealth{Name: "address"}
		available := m.address.IsAddressAvailable(ctx)
		if !available {
			health.Status = "unhealthy"
			health.Error = "address search not available"
		} else {
			health.Status = "healthy"
		}
		result.Services["address"] = health
	}

	// Set overall status
	result.OverallStatus = "healthy"
	for _, service := range result.Services {
		if service.Status != "healthy" {
			result.OverallStatus = "degraded"
			break
		}
	}

	return result
}

// HealthCheckResult represents the result of external API health checks
type HealthCheckResult struct {
	OverallStatus string                    `json:"overall_status"`
	Services      map[string]*ServiceHealth `json:"services"`
}

// ServiceHealth represents the health status of a single external service
type ServiceHealth struct {
	Name   string `json:"name"`
	Status string `json:"status"` // "healthy", "unhealthy"
	Error  string `json:"error,omitempty"`
}

// IsHealthy returns true if all services are healthy
func (r *HealthCheckResult) IsHealthy() bool {
	return r.OverallStatus == "healthy"
}

// GetUnhealthyServices returns a list of unhealthy services
func (r *HealthCheckResult) GetUnhealthyServices() []string {
	var unhealthy []string
	for _, service := range r.Services {
		if service.Status != "healthy" {
			unhealthy = append(unhealthy, service.Name)
		}
	}
	return unhealthy
}