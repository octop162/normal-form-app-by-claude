// Package dto defines data transfer objects for option management.
package dto

// OptionResponse represents an option in API responses
type OptionResponse struct {
	ID                int    `json:"id"`
	OptionType        string `json:"option_type"`
	OptionName        string `json:"option_name"`
	Description       string `json:"description,omitempty"`
	PlanCompatibility string `json:"plan_compatibility"`
	IsActive          bool   `json:"is_active"`
}

// OptionsGetRequest represents the request for getting available options
type OptionsGetRequest struct {
	PlanType string `form:"plan_type" validate:"required,oneof=A B"`
	Region   string `form:"region" validate:"omitempty"`
}

// OptionsGetResponse represents the response for getting available options
type OptionsGetResponse struct {
	Options []OptionResponse `json:"options"`
}

// InventoryCheckRequest represents the request for checking option inventory
type InventoryCheckRequest struct {
	OptionTypes []string `json:"option_types" validate:"required,dive,oneof=AA BB AB"`
}

// InventoryCheckResponse represents the response for inventory check
type InventoryCheckResponse struct {
	Inventory map[string]int `json:"inventory"`
}
