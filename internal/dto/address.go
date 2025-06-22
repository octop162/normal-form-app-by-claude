// Package dto defines data transfer objects for address management.
package dto

// AddressSearchRequest represents the request for address search
type AddressSearchRequest struct {
	PostalCode string `form:"postal_code" validate:"required,len=7,numeric"`
}

// AddressSearchResponse represents the response for address search
type AddressSearchResponse struct {
	Found      bool   `json:"found"`
	Prefecture string `json:"prefecture,omitempty"`
	City       string `json:"city,omitempty"`
	Town       string `json:"town,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}

// RegionCheckRequest represents the request for region restriction check
type RegionCheckRequest struct {
	Prefecture  string   `json:"prefecture" validate:"required"`
	City        string   `json:"city" validate:"required"`
	OptionTypes []string `json:"option_types" validate:"required,dive,oneof=AA BB AB"`
}

// RegionCheckResponse represents the response for region restriction check
type RegionCheckResponse struct {
	Restrictions map[string]bool `json:"restrictions"`
}

// PrefectureResponse represents a prefecture in API responses
type PrefectureResponse struct {
	ID             int    `json:"id"`
	PrefectureCode string `json:"prefecture_code"`
	PrefectureName string `json:"prefecture_name"`
	Region         string `json:"region"`
}

// PrefecturesGetResponse represents the response for getting prefectures
type PrefecturesGetResponse struct {
	Prefectures []PrefectureResponse `json:"prefectures"`
}
