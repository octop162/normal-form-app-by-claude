// Package handler provides error handling utilities.
package handler

import (
	"strings"
)

// isValidationError checks if the error is a validation error
func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	validationKeywords := []string{
		"validation",
		"invalid",
		"required",
		"format",
		"length",
	}

	for _, keyword := range validationKeywords {
		if strings.Contains(strings.ToLower(errMsg), keyword) {
			return true
		}
	}

	return false
}

// isDuplicateError checks if the error is a duplicate/conflict error
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	duplicateKeywords := []string{
		"already exists",
		"duplicate",
		"conflict",
		"unique constraint",
	}

	for _, keyword := range duplicateKeywords {
		if strings.Contains(strings.ToLower(errMsg), keyword) {
			return true
		}
	}

	return false
}

// isNotFoundError checks if the error is a not found error
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	notFoundKeywords := []string{
		"not found",
		"does not exist",
		"no rows",
	}

	for _, keyword := range notFoundKeywords {
		if strings.Contains(strings.ToLower(errMsg), keyword) {
			return true
		}
	}

	return false
}

// isExpiredError checks if the error is related to expiration
func isExpiredError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	expiredKeywords := []string{
		"expired",
		"timeout",
		"timed out",
	}

	for _, keyword := range expiredKeywords {
		if strings.Contains(strings.ToLower(errMsg), keyword) {
			return true
		}
	}

	return false
}
