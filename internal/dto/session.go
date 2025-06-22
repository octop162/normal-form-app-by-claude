// Package dto defines data transfer objects for session management.
package dto

import (
	"time"
)

// SessionCreateRequest represents the request for creating a session
type SessionCreateRequest struct {
	UserData map[string]interface{} `json:"user_data" validate:"required"`
}

// SessionCreateResponse represents the response for session creation
type SessionCreateResponse struct {
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionUpdateRequest represents the request for updating a session
type SessionUpdateRequest struct {
	UserData map[string]interface{} `json:"user_data" validate:"required"`
}

// SessionUpdateResponse represents the response for session update
type SessionUpdateResponse struct {
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SessionGetResponse represents the response for session retrieval
type SessionGetResponse struct {
	SessionID string                 `json:"session_id"`
	UserData  map[string]interface{} `json:"user_data"`
	ExpiresAt time.Time              `json:"expires_at"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// SessionDeleteResponse represents the response for session deletion
type SessionDeleteResponse struct {
	Message string `json:"message"`
}
