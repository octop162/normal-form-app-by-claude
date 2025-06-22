// Package model defines domain models for the application.
package model

import (
	"time"
)

// User represents a registered user
type User struct {
	ID           int       `json:"id" db:"id"`
	LastName     string    `json:"last_name" db:"last_name"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastNameKana string    `json:"last_name_kana" db:"last_name_kana"`
	FirstNameKana string   `json:"first_name_kana" db:"first_name_kana"`
	Phone1       string    `json:"phone1" db:"phone1"`
	Phone2       string    `json:"phone2" db:"phone2"`
	Phone3       string    `json:"phone3" db:"phone3"`
	PostalCode1  string    `json:"postal_code1" db:"postal_code1"`
	PostalCode2  string    `json:"postal_code2" db:"postal_code2"`
	Prefecture   string    `json:"prefecture" db:"prefecture"`
	City         string    `json:"city" db:"city"`
	Town         *string   `json:"town" db:"town"`
	Chome        *string   `json:"chome" db:"chome"`
	Banchi       string    `json:"banchi" db:"banchi"`
	Go           *string   `json:"go" db:"go"`
	Building     *string   `json:"building" db:"building"`
	Room         *string   `json:"room" db:"room"`
	Email        string    `json:"email" db:"email"`
	PlanType     string    `json:"plan_type" db:"plan_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserOption represents a selected option for a user
type UserOption struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	OptionType string    `json:"option_type" db:"option_type"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// UserSession represents a temporary session for form data
type UserSession struct {
	ID        string                 `json:"id" db:"id"`
	UserData  map[string]interface{} `json:"user_data" db:"user_data"`
	ExpiresAt time.Time              `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" db:"updated_at"`
}

// OptionMaster represents master data for options
type OptionMaster struct {
	ID                int       `json:"id" db:"id"`
	OptionType        string    `json:"option_type" db:"option_type"`
	OptionName        string    `json:"option_name" db:"option_name"`
	Description       *string   `json:"description" db:"description"`
	PlanCompatibility string    `json:"plan_compatibility" db:"plan_compatibility"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// PrefectureMaster represents master data for prefectures
type PrefectureMaster struct {
	ID             int       `json:"id" db:"id"`
	PrefectureCode string    `json:"prefecture_code" db:"prefecture_code"`
	PrefectureName string    `json:"prefecture_name" db:"prefecture_name"`
	Region         string    `json:"region" db:"region"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.LastName + " " + u.FirstName
}

// GetFullNameKana returns the full name in katakana
func (u *User) GetFullNameKana() string {
	return u.LastNameKana + " " + u.FirstNameKana
}

// GetPhoneNumber returns the complete phone number
func (u *User) GetPhoneNumber() string {
	return u.Phone1 + "-" + u.Phone2 + "-" + u.Phone3
}

// GetPostalCode returns the complete postal code
func (u *User) GetPostalCode() string {
	return u.PostalCode1 + "-" + u.PostalCode2
}

// GetFullAddress returns the complete address
func (u *User) GetFullAddress() string {
	address := u.Prefecture + u.City
	
	if u.Town != nil && *u.Town != "" {
		address += *u.Town
	}
	
	if u.Chome != nil && *u.Chome != "" {
		address += *u.Chome
	}
	
	address += u.Banchi
	
	if u.Go != nil && *u.Go != "" {
		address += "-" + *u.Go
	}
	
	if u.Building != nil && *u.Building != "" {
		address += " " + *u.Building
	}
	
	if u.Room != nil && *u.Room != "" {
		address += " " + *u.Room
	}
	
	return address
}

// IsExpired checks if the session is expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// CanUseOption checks if the option is compatible with the user's plan
func (u *User) CanUseOption(option *OptionMaster) bool {
	if !option.IsActive {
		return false
	}
	
	switch option.PlanCompatibility {
	case "A":
		return u.PlanType == "A"
	case "B":
		return u.PlanType == "B"
	case "AB":
		return u.PlanType == "A" || u.PlanType == "B"
	default:
		return false
	}
}

// Address represents address information for external APIs
type Address struct {
	PostalCode string `json:"postal_code"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Town       string `json:"town,omitempty"`
}

// Plan represents plan information
type Plan struct {
	PlanType    string `json:"plan_type"`
	PlanName    string `json:"plan_name"`
	Description string `json:"description,omitempty"`
}