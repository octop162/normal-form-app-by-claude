// Package service provides address management business logic.
package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/internal/repository"
	"github.com/octop162/normal-form-app-by-claude/pkg/external"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

const (
	// Postal code validation constants
	postalCodeLength = 7
)

// AddressService defines the interface for address business logic
type AddressService interface {
	SearchByPostalCode(ctx context.Context, req *dto.AddressSearchRequest) (*dto.AddressSearchResponse, error)
	CheckRegionRestrictions(ctx context.Context, req *dto.RegionCheckRequest) (*dto.RegionCheckResponse, error)
	GetPrefectures(ctx context.Context) (*dto.PrefecturesGetResponse, error)
	GetPrefectureByName(ctx context.Context, name string) (*dto.PrefectureResponse, error)
}

// addressService implements AddressService
type addressService struct {
	prefectureRepo repository.PrefectureRepository
	externalAPI    *external.Manager
	log            *logger.Logger
}

// NewAddressService creates a new address service
func NewAddressService(
	prefectureRepo repository.PrefectureRepository,
	externalAPI *external.Manager,
	log *logger.Logger,
) AddressService {
	return &addressService{
		prefectureRepo: prefectureRepo,
		externalAPI:    externalAPI,
		log:            log,
	}
}

// SearchByPostalCode searches for address information by postal code
func (s *addressService) SearchByPostalCode(
	ctx context.Context, req *dto.AddressSearchRequest,
) (*dto.AddressSearchResponse, error) {
	// Validate postal code format (should be 7 digits)
	if len(req.PostalCode) != postalCodeLength {
		return &dto.AddressSearchResponse{
			Found: false,
		}, nil
	}

	// Try external address API first if available
	if s.externalAPI != nil && s.externalAPI.AddressClient() != nil {
		addressInfo, err := s.externalAPI.AddressClient().SearchByPostalCode(ctx, req.PostalCode)
		if err != nil {
			s.log.WithError(err).WithField("postal_code", req.PostalCode).Warn("External address API failed, falling back to mock data")
		} else {
			return &dto.AddressSearchResponse{
				Found:      true,
				Prefecture: addressInfo.Prefecture,
				City:       addressInfo.City,
				Town:       addressInfo.Town,
				PostalCode: formatPostalCode(req.PostalCode),
			}, nil
		}
	}

	// Fallback to mock data
	address := s.getMockAddressData(req.PostalCode)
	if address == nil {
		return &dto.AddressSearchResponse{
			Found: false,
		}, nil
	}

	return &dto.AddressSearchResponse{
		Found:      true,
		Prefecture: address.Prefecture,
		City:       address.City,
		Town:       address.Town,
		PostalCode: formatPostalCode(req.PostalCode),
	}, nil
}

// CheckRegionRestrictions checks if options are available in the specified region
func (s *addressService) CheckRegionRestrictions(
	ctx context.Context, req *dto.RegionCheckRequest,
) (*dto.RegionCheckResponse, error) {
	restrictions := make(map[string]bool)

	// Try external region API first if available
	if s.externalAPI != nil && s.externalAPI.RegionClient() != nil {
		regionRestrictions, err := s.externalAPI.RegionClient().CheckRegionRestrictions(
			ctx, req.Prefecture, req.City, req.OptionTypes,
		)
		if err != nil {
			s.log.WithError(err).
				WithField("prefecture", req.Prefecture).
				WithField("city", req.City).
				WithField("options", req.OptionTypes).
				Warn("External region API failed, falling back to local logic")
		} else {
			return &dto.RegionCheckResponse{
				Restrictions: regionRestrictions,
			}, nil
		}
	}

	// Fallback to local logic
	prefecture, err := s.prefectureRepo.GetByName(ctx, req.Prefecture)
	if err != nil {
		s.log.WithError(err).WithField("prefecture", req.Prefecture).Error("Failed to get prefecture")
		return nil, fmt.Errorf("failed to get prefecture: %w", err)
	}

	// Check restrictions for each option type using local logic
	for _, optionType := range req.OptionTypes {
		allowed := s.checkOptionAllowedInRegion(prefecture, req.City, optionType)
		restrictions[optionType] = allowed
	}

	return &dto.RegionCheckResponse{
		Restrictions: restrictions,
	}, nil
}

// GetPrefectures retrieves all active prefectures
func (s *addressService) GetPrefectures(ctx context.Context) (*dto.PrefecturesGetResponse, error) {
	prefectures, err := s.prefectureRepo.GetActive(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to get prefectures")
		return nil, fmt.Errorf("failed to get prefectures: %w", err)
	}

	// Convert to response DTOs
	prefectureResponses := make([]dto.PrefectureResponse, len(prefectures))
	for i, prefecture := range prefectures {
		prefectureResponses[i] = s.convertPrefectureToResponse(prefecture)
	}

	return &dto.PrefecturesGetResponse{
		Prefectures: prefectureResponses,
	}, nil
}

// GetPrefectureByName retrieves a specific prefecture by name
func (s *addressService) GetPrefectureByName(ctx context.Context, name string) (*dto.PrefectureResponse, error) {
	prefecture, err := s.prefectureRepo.GetByName(ctx, name)
	if err != nil {
		s.log.WithError(err).WithField("prefecture_name", name).Error("Failed to get prefecture by name")
		return nil, fmt.Errorf("failed to get prefecture by name: %w", err)
	}

	response := s.convertPrefectureToResponse(prefecture)
	return &response, nil
}

// getMockAddressData returns mock address data for testing
// TODO: Replace with actual external postal code API call
func (s *addressService) getMockAddressData(postalCode string) *model.Address {
	// Mock address data for common postal codes
	mockData := map[string]*model.Address{
		"1000001": {
			PostalCode: "100-0001",
			Prefecture: "東京都",
			City:       "千代田区",
			Town:       "千代田",
		},
		"1500002": {
			PostalCode: "150-0002",
			Prefecture: "東京都",
			City:       "渋谷区",
			Town:       "渋谷",
		},
		"5410041": {
			PostalCode: "541-0041",
			Prefecture: "大阪府",
			City:       "大阪市中央区",
			Town:       "北浜",
		},
		"2310023": {
			PostalCode: "231-0023",
			Prefecture: "神奈川県",
			City:       "横浜市中区",
			Town:       "山下町",
		},
		"4600008": {
			PostalCode: "460-0008",
			Prefecture: "愛知県",
			City:       "名古屋市中区",
			Town:       "栄",
		},
	}

	return mockData[postalCode]
}

// checkOptionAllowedInRegion checks if an option is allowed in the specified region
// TODO: Implement actual region restriction logic
func (s *addressService) checkOptionAllowedInRegion(
	prefecture *model.PrefectureMaster, _ string, optionType string,
) bool {
	// Mock region restrictions for testing
	// In production, this would call external region restriction API

	// Example restrictions:
	// - AA option not available in certain remote areas
	// - BB option restricted in some metropolitan areas
	switch optionType {
	case "AA":
		// AA option not available in Hokkaido for this example
		return prefecture.PrefectureName != "北海道"
	case "BB":
		// BB option not available in major metropolitan areas
		restrictedCities := []string{"東京都", "大阪府", "愛知県"}
		return !slices.Contains(restrictedCities, prefecture.PrefectureName)
	case "AB":
		// AB option available everywhere
		return true
	default:
		return false
	}
}

// convertPrefectureToResponse converts prefecture model to response DTO
func (s *addressService) convertPrefectureToResponse(prefecture *model.PrefectureMaster) dto.PrefectureResponse {
	return dto.PrefectureResponse{
		ID:             prefecture.ID,
		PrefectureCode: prefecture.PrefectureCode,
		PrefectureName: prefecture.PrefectureName,
		Region:         prefecture.Region,
	}
}

// formatPostalCode formats postal code with hyphen (XXXXXXX -> XXX-XXXX)
func formatPostalCode(postalCode string) string {
	if len(postalCode) != postalCodeLength {
		return postalCode
	}
	return postalCode[:3] + "-" + postalCode[3:]
}
