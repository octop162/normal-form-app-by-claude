// Package service provides plan management business logic.
package service

import (
	"context"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// PlanService defines the interface for plan business logic
type PlanService interface {
	GetAvailablePlans(ctx context.Context) (*dto.PlansGetResponse, error)
	GetPlanByType(ctx context.Context, planType string) (*dto.PlanResponse, error)
	ValidatePlanType(ctx context.Context, planType string) (bool, error)
}

// planService implements PlanService
type planService struct {
	log *logger.Logger
}

// NewPlanService creates a new plan service
func NewPlanService(log *logger.Logger) PlanService {
	return &planService{
		log: log,
	}
}

// GetAvailablePlans retrieves all available plans
func (s *planService) GetAvailablePlans(_ context.Context) (*dto.PlansGetResponse, error) {
	// TODO: In production, this might come from a database or external service
	// For now, return static plan data
	plans := []dto.PlanResponse{
		{
			PlanType:    "A",
			PlanName:    "Aプラン",
			Description: "基本プランです。標準的なサービスをご利用いただけます。",
		},
		{
			PlanType:    "B",
			PlanName:    "Bプラン",
			Description: "プレミアムプランです。より充実したサービスをご利用いただけます。",
		},
	}

	return &dto.PlansGetResponse{
		Plans: plans,
	}, nil
}

// GetPlanByType retrieves a specific plan by type
func (s *planService) GetPlanByType(ctx context.Context, planType string) (*dto.PlanResponse, error) {
	plans, err := s.GetAvailablePlans(ctx)
	if err != nil {
		return nil, err
	}

	for _, plan := range plans.Plans {
		if plan.PlanType == planType {
			return &plan, nil
		}
	}

	return nil, fmt.Errorf("plan type %s not found", planType)
}

// ValidatePlanType validates if a plan type is valid
func (s *planService) ValidatePlanType(ctx context.Context, planType string) (bool, error) {
	_, err := s.GetPlanByType(ctx, planType)
	if err != nil {
		return false, nil // Plan type not found, but no error
	}

	return true, nil
}
