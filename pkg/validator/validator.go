// Package validator provides validation functionality for the application.
package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const (
	// Phone number validation constants
	freeDial4DigitLength = 4
	mobileNumberLength   = 11
	freeDial3DigitLength = 3
)

var (
	// Katakana regex pattern
	katakanaPattern = regexp.MustCompile(`^[ァ-ヶー]+$`)
	// Numeric regex pattern
	numericPattern = regexp.MustCompile(`^[0-9]+$`)
)

// CustomValidator wraps the validator with custom validation rules
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator instance with custom rules
func NewValidator() (*CustomValidator, error) {
	v := validator.New()

	// Register custom validation functions
	if err := v.RegisterValidation("katakana", validateKatakana); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("numeric", validateNumeric); err != nil {
		return nil, err
	}
	if err := v.RegisterValidation("phone", validatePhone); err != nil {
		return nil, err
	}

	return &CustomValidator{validator: v}, nil
}

// ValidateStruct validates a struct using the configured validator
func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	return cv.validator.Struct(s)
}

// GetValidator returns the underlying validator instance
func (cv *CustomValidator) GetValidator() *validator.Validate {
	return cv.validator
}

// validateKatakana validates that the field contains only katakana characters
func validateKatakana(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty values are handled by required tag
	}
	return katakanaPattern.MatchString(value)
}

// validateNumeric validates that the field contains only numeric characters
func validateNumeric(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty values are handled by required tag
	}
	return numericPattern.MatchString(value)
}

// validatePhone validates phone number format and restrictions
func validatePhone(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // Empty values are handled by required tag
	}

	// Check for forbidden numbers (free dial numbers)
	if len(value) >= freeDial4DigitLength {
		prefix := value[:freeDial4DigitLength]
		if prefix == "0120" || prefix == "0800" {
			return false
		}
	}

	// For 11-digit numbers, check if it starts with 0X0 (mobile numbers)
	if len(value) == mobileNumberLength {
		if len(value) >= freeDial3DigitLength && value[0] == '0' && value[2] == '0' {
			secondDigit := value[1]
			// Valid mobile number prefixes: 070, 080, 090
			return secondDigit == '7' || secondDigit == '8' || secondDigit == '9'
		}
		return false
	}

	return true
}

// IsValidEmail performs basic email validation
func IsValidEmail(email string) bool {
	// Basic email regex - more comprehensive validation can be added
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(email)
}

// IsValidPostalCode validates Japanese postal code format (XXX-XXXX)
func IsValidPostalCode(postalCode string) bool {
	postalPattern := regexp.MustCompile(`^[0-9]{3}-[0-9]{4}$`)
	return postalPattern.MatchString(postalCode)
}

// IsValidPlanType validates plan type
func IsValidPlanType(planType string) bool {
	return planType == "A" || planType == "B"
}

// IsValidOptionType validates option type
func IsValidOptionType(optionType string) bool {
	return optionType == "AA" || optionType == "BB" || optionType == "AB"
}

// IsValidPhone validates phone number with business rules
func IsValidPhone(phoneNumber string) bool {
	if phoneNumber == "" {
		return false
	}

	// Check for forbidden numbers (free dial numbers)
	if len(phoneNumber) >= freeDial4DigitLength {
		prefix := phoneNumber[:freeDial4DigitLength]
		if prefix == "0120" || prefix == "0800" {
			return false
		}
	}

	// For 11-digit numbers, check if it starts with 0X0 (mobile numbers)
	if len(phoneNumber) == mobileNumberLength {
		if len(phoneNumber) >= freeDial3DigitLength && phoneNumber[0] == '0' && phoneNumber[2] == '0' {
			secondDigit := phoneNumber[1]
			// Valid mobile number prefixes: 070, 080, 090
			return secondDigit == '7' || secondDigit == '8' || secondDigit == '9'
		}
		return false
	}

	// For other lengths, basic numeric validation
	return numericPattern.MatchString(phoneNumber)
}

// ContainsOnlyKatakana checks if string contains only katakana characters
func ContainsOnlyKatakana(s string) bool {
	for _, r := range s {
		if !unicode.In(r, unicode.Katakana) && r != 'ー' {
			return false
		}
	}
	return true
}
