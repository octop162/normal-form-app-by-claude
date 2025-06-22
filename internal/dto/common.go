// Package dto defines common data transfer objects for API communication.
package dto

// APIResponse represents the standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error in API responses
type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// PingResponse represents the response for ping endpoint
type PingResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// HealthResponse represents the response for health check
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// SimpleStatusResponse represents a simple status response
type SimpleStatusResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// PlansGetResponse represents the response for getting available plans
type PlansGetResponse struct {
	Plans []PlanResponse `json:"plans"`
}

// PlanResponse represents a plan in API responses
type PlanResponse struct {
	PlanType    string `json:"plan_type"`
	PlanName    string `json:"plan_name"`
	Description string `json:"description,omitempty"`
}
