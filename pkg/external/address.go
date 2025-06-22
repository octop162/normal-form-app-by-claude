// Package external provides address search API client functionality.
package external

import (
	"context"
	"fmt"
	"regexp"

	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	addressSearchEndpoint = "/api/address/search"
)

var (
	// Postal code validation regex (3 digits + 4 digits)
	postalCodeRegex = regexp.MustCompile(`^\d{3}-?\d{4}$`)
)

// AddressClient handles address search-related external API calls
type AddressClient struct {
	client *Client
	log    *logger.Logger
}

// NewAddressClient creates a new address API client
func NewAddressClient(config *Config, log *logger.Logger) *AddressClient {
	return &AddressClient{
		client: NewClient(config, log),
		log:    log,
	}
}

// AddressSearchRequest represents the request payload for address search
type AddressSearchRequest struct {
	PostalCode string `json:"postal_code" validate:"required"`
}

// AddressSearchResponse represents the response from address search API
type AddressSearchResponse struct {
	Success bool         `json:"success"`
	Data    *AddressData `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// AddressData represents the address information returned by the API
type AddressData struct {
	PostalCode string `json:"postal_code"`
	Prefecture string `json:"prefecture"`
	City       string `json:"city"`
	Town       string `json:"town,omitempty"`
}

// AddressInfo represents complete address information
type AddressInfo struct {
	PostalCode1 string `json:"postal_code_1"` // First 3 digits
	PostalCode2 string `json:"postal_code_2"` // Last 4 digits
	Prefecture  string `json:"prefecture"`
	City        string `json:"city"`
	Town        string `json:"town,omitempty"`
	FullAddress string `json:"full_address"`
}

// SearchByPostalCode searches for address information using postal code
func (ac *AddressClient) SearchByPostalCode(ctx context.Context, postalCode string) (*AddressInfo, error) {
	if postalCode == "" {
		return nil, fmt.Errorf("postal code cannot be empty")
	}

	// Validate postal code format
	if !postalCodeRegex.MatchString(postalCode) {
		return nil, fmt.Errorf("invalid postal code format: %s", postalCode)
	}

	// Normalize postal code (remove hyphen if present)
	normalizedPostalCode := normalizePostalCode(postalCode)

	// Prepare request
	req := &AddressSearchRequest{
		PostalCode: normalizedPostalCode,
	}

	// Make API call
	var resp AddressSearchResponse
	err := ac.client.PostJSON(ctx, addressSearchEndpoint, req, &resp)
	if err != nil {
		ac.log.WithError(err).WithField("postal_code", postalCode).Error("Failed to search address")
		return nil, fmt.Errorf("address search API call failed: %w", err)
	}

	// Validate response
	if !resp.Success {
		errMsg := "unknown error"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		ac.log.WithField("postal_code", postalCode).WithField("api_error", errMsg).Error("Address API returned error")
		return nil, fmt.Errorf("address API error: %s", errMsg)
	}

	if resp.Data == nil {
		ac.log.WithField("postal_code", postalCode).Error("Address API returned no data")
		return nil, fmt.Errorf("no address data found for postal code: %s", postalCode)
	}

	// Convert to AddressInfo format
	addressInfo := &AddressInfo{
		PostalCode1: normalizedPostalCode[:3],
		PostalCode2: normalizedPostalCode[3:],
		Prefecture:  resp.Data.Prefecture,
		City:        resp.Data.City,
		Town:        resp.Data.Town,
		FullAddress: buildFullAddress(resp.Data),
	}

	ac.log.WithField("postal_code", postalCode).WithField("address_info", addressInfo).Debug("Address search completed")
	return addressInfo, nil
}

// SearchByPostalCodeParts searches for address information using postal code parts
func (ac *AddressClient) SearchByPostalCodeParts(ctx context.Context, postalCode1, postalCode2 string) (*AddressInfo, error) {
	if postalCode1 == "" || postalCode2 == "" {
		return nil, fmt.Errorf("postal code parts cannot be empty")
	}

	if len(postalCode1) != 3 || len(postalCode2) != 4 {
		return nil, fmt.Errorf("invalid postal code parts format: %s-%s", postalCode1, postalCode2)
	}

	fullPostalCode := postalCode1 + postalCode2
	return ac.SearchByPostalCode(ctx, fullPostalCode)
}

// ValidatePostalCode validates the format of a postal code
func (ac *AddressClient) ValidatePostalCode(postalCode string) error {
	if postalCode == "" {
		return fmt.Errorf("postal code cannot be empty")
	}

	if !postalCodeRegex.MatchString(postalCode) {
		return fmt.Errorf("invalid postal code format: %s (expected format: XXX-XXXX or XXXXXXX)", postalCode)
	}

	return nil
}

// ValidatePostalCodeParts validates the format of postal code parts
func (ac *AddressClient) ValidatePostalCodeParts(postalCode1, postalCode2 string) error {
	if postalCode1 == "" || postalCode2 == "" {
		return fmt.Errorf("postal code parts cannot be empty")
	}

	if len(postalCode1) != 3 {
		return fmt.Errorf("postal code first part must be 3 digits: %s", postalCode1)
	}

	if len(postalCode2) != 4 {
		return fmt.Errorf("postal code second part must be 4 digits: %s", postalCode2)
	}

	// Check if all characters are digits
	digitRegex := regexp.MustCompile(`^\d+$`)
	if !digitRegex.MatchString(postalCode1) {
		return fmt.Errorf("postal code first part must contain only digits: %s", postalCode1)
	}

	if !digitRegex.MatchString(postalCode2) {
		return fmt.Errorf("postal code second part must contain only digits: %s", postalCode2)
	}

	return nil
}

// normalizePostalCode removes hyphen from postal code
func normalizePostalCode(postalCode string) string {
	if len(postalCode) == 8 && postalCode[3] == '-' {
		return postalCode[:3] + postalCode[4:]
	}
	return postalCode
}

// buildFullAddress constructs a full address string from address data
func buildFullAddress(data *AddressData) string {
	address := data.Prefecture + data.City
	if data.Town != "" {
		address += data.Town
	}
	return address
}

// IsAddressAvailable checks if address search is available (API health check)
func (ac *AddressClient) IsAddressAvailable(ctx context.Context) bool {
	// Try searching with a known valid postal code (Tokyo Station)
	_, err := ac.SearchByPostalCode(ctx, "1000005")
	return err == nil
}