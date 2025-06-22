// Package handler provides constants for HTTP handlers.
package handler

// HTTP Error Codes
const (
	// Generic errors
	ErrorCodeInvalidRequest  = "INVALID_REQUEST"
	ErrorCodeInternalError   = "INTERNAL_ERROR"
	ErrorCodeValidationError = "VALIDATION_ERROR"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeDuplicateError  = "DUPLICATE_ERROR"

	// User-specific errors
	ErrorCodeUserNotFound  = "USER_NOT_FOUND"
	ErrorCodeInvalidUserID = "INVALID_USER_ID"

	// Session-specific errors
	ErrorCodeSessionNotFound     = "SESSION_NOT_FOUND"
	ErrorCodeSessionCreateFailed = "SESSION_CREATE_FAILED"
	ErrorCodeMissingSessionID    = "MISSING_SESSION_ID"

	// Option-specific errors
	ErrorCodeOptionNotFound       = "OPTION_NOT_FOUND"
	ErrorCodeMissingOptionType    = "MISSING_OPTION_TYPE"
	ErrorCodeInventoryCheckFailed = "INVENTORY_CHECK_FAILED"

	// Address-specific errors
	ErrorCodeAddressSearchFailed   = "ADDRESS_SEARCH_FAILED"
	ErrorCodeRegionCheckFailed     = "REGION_CHECK_FAILED"
	ErrorCodePrefectureNotFound    = "PREFECTURE_NOT_FOUND"
	ErrorCodeMissingPrefectureName = "MISSING_PREFECTURE_NAME"

	// Plan-specific errors
	ErrorCodePlanNotFound    = "PLAN_NOT_FOUND"
	ErrorCodeMissingPlanType = "MISSING_PLAN_TYPE"
)

// HTTP Error Messages
const (
	MessageInvalidRequest     = "Invalid request format"
	MessageInvalidQueryParams = "Invalid query parameters"
	MessageInternalError      = "Internal server error"
	MessageValidationFailed   = "Validation failed"
	MessageUserNotFound       = "User not found"
	MessageSessionNotFound    = "Session not found or expired"
	MessageOptionNotFound     = "Option not found"
	MessagePrefectureNotFound = "Prefecture not found"
	MessagePlanNotFound       = "Plan not found"
)
