// Package dto defines data transfer objects for API communication.
package dto

import (
	"time"
)

// UserCreateRequest represents the request for user registration
type UserCreateRequest struct {
	LastName      string   `json:"last_name" validate:"required,max=15"`
	FirstName     string   `json:"first_name" validate:"required,max=15"`
	LastNameKana  string   `json:"last_name_kana" validate:"required,max=15,katakana"`
	FirstNameKana string   `json:"first_name_kana" validate:"required,max=15,katakana"`
	Phone1        string   `json:"phone1" validate:"required,len=3,numeric"`
	Phone2        string   `json:"phone2" validate:"required,min=1,max=4,numeric"`
	Phone3        string   `json:"phone3" validate:"required,len=4,numeric"`
	PostalCode1   string   `json:"postal_code1" validate:"required,len=3,numeric"`
	PostalCode2   string   `json:"postal_code2" validate:"required,len=4,numeric"`
	Prefecture    string   `json:"prefecture" validate:"required,max=10"`
	City          string   `json:"city" validate:"required,max=50"`
	Town          *string  `json:"town" validate:"omitempty,max=50"`
	Chome         *string  `json:"chome" validate:"omitempty,max=10"`
	Banchi        string   `json:"banchi" validate:"required,max=10"`
	Go            *string  `json:"go" validate:"omitempty,max=10"`
	Building      *string  `json:"building" validate:"omitempty,max=100"`
	Room          *string  `json:"room" validate:"omitempty,max=20"`
	Email         string   `json:"email" validate:"required,email,max=256"`
	EmailConfirm  string   `json:"email_confirm" validate:"required,eqfield=Email"`
	PlanType      string   `json:"plan_type" validate:"required,oneof=A B"`
	OptionTypes   []string `json:"option_types" validate:"dive,oneof=AA BB AB"`
}

// UserCreateResponse represents the response for user registration
type UserCreateResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

// UserValidateRequest represents the request for user data validation
type UserValidateRequest struct {
	UserCreateRequest
}

// UserValidateResponse represents the response for user data validation
type UserValidateResponse struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

// UserResponse represents a user in API responses
type UserResponse struct {
	ID            int       `json:"id"`
	LastName      string    `json:"last_name"`
	FirstName     string    `json:"first_name"`
	LastNameKana  string    `json:"last_name_kana"`
	FirstNameKana string    `json:"first_name_kana"`
	PhoneNumber   string    `json:"phone_number"`
	PostalCode    string    `json:"postal_code"`
	Address       string    `json:"address"`
	Email         string    `json:"email"`
	PlanType      string    `json:"plan_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
