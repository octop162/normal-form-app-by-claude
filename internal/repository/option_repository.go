// Package repository provides option master data access functionality.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// OptionRepository defines the interface for option master data access
type OptionRepository interface {
	GetAll(ctx context.Context) ([]*model.OptionMaster, error)
	GetByPlanType(ctx context.Context, planType string) ([]*model.OptionMaster, error)
	GetByOptionType(ctx context.Context, optionType string) (*model.OptionMaster, error)
	GetActiveOptions(ctx context.Context) ([]*model.OptionMaster, error)
	GetCompatibleOptions(ctx context.Context, planType string) ([]*model.OptionMaster, error)
}

// optionRepository implements OptionRepository
type optionRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewOptionRepository creates a new option repository
func NewOptionRepository(db *sql.DB, log *logger.Logger) OptionRepository {
	return &optionRepository{
		db:  db,
		log: log,
	}
}

// GetAll retrieves all option master data
func (r *optionRepository) GetAll(ctx context.Context) ([]*model.OptionMaster, error) {
	query := `
		SELECT id, option_type, option_name, description, plan_compatibility, is_active, created_at, updated_at
		FROM options_master
		ORDER BY option_type ASC`

	return r.queryOptions(ctx, query)
}

// GetByPlanType retrieves options compatible with a specific plan type
func (r *optionRepository) GetByPlanType(ctx context.Context, planType string) ([]*model.OptionMaster, error) {
	query := `
		SELECT id, option_type, option_name, description, plan_compatibility, is_active, created_at, updated_at
		FROM options_master
		WHERE is_active = true AND (plan_compatibility = $1 OR plan_compatibility = 'AB')
		ORDER BY option_type ASC`

	rows, err := r.db.QueryContext(ctx, query, planType)
	if err != nil {
		r.log.WithError(err).WithField("plan_type", planType).Error("Failed to get options by plan type")
		return nil, fmt.Errorf("failed to get options by plan type: %w", err)
	}
	defer rows.Close()

	return r.scanOptions(rows)
}

// GetByOptionType retrieves a specific option by option type
func (r *optionRepository) GetByOptionType(ctx context.Context, optionType string) (*model.OptionMaster, error) {
	query := `
		SELECT id, option_type, option_name, description, plan_compatibility, is_active, created_at, updated_at
		FROM options_master
		WHERE option_type = $1`

	var option model.OptionMaster
	err := r.db.QueryRowContext(ctx, query, optionType).Scan(
		&option.ID, &option.OptionType, &option.OptionName, &option.Description,
		&option.PlanCompatibility, &option.IsActive, &option.CreatedAt, &option.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("option not found: %w", err)
		}
		r.log.WithError(err).WithField("option_type", optionType).Error("Failed to get option by type")
		return nil, fmt.Errorf("failed to get option by type: %w", err)
	}

	return &option, nil
}

// GetActiveOptions retrieves all active options
func (r *optionRepository) GetActiveOptions(ctx context.Context) ([]*model.OptionMaster, error) {
	query := `
		SELECT id, option_type, option_name, description, plan_compatibility, is_active, created_at, updated_at
		FROM options_master
		WHERE is_active = true
		ORDER BY option_type ASC`

	return r.queryOptions(ctx, query)
}

// GetCompatibleOptions retrieves options compatible with a specific plan type (active only)
func (r *optionRepository) GetCompatibleOptions(ctx context.Context, planType string) ([]*model.OptionMaster, error) {
	return r.GetByPlanType(ctx, planType)
}

// queryOptions executes a query and returns options
func (r *optionRepository) queryOptions(
	ctx context.Context, query string, args ...any,
) ([]*model.OptionMaster, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.WithError(err).Error("Failed to query options")
		return nil, fmt.Errorf("failed to query options: %w", err)
	}
	defer rows.Close()

	return r.scanOptions(rows)
}

// scanOptions scans rows into option master structs
func (r *optionRepository) scanOptions(rows *sql.Rows) ([]*model.OptionMaster, error) {
	var options []*model.OptionMaster

	for rows.Next() {
		var option model.OptionMaster
		err := rows.Scan(
			&option.ID, &option.OptionType, &option.OptionName, &option.Description,
			&option.PlanCompatibility, &option.IsActive, &option.CreatedAt, &option.UpdatedAt,
		)
		if err != nil {
			r.log.WithError(err).Error("Failed to scan option row")
			return nil, fmt.Errorf("failed to scan option row: %w", err)
		}
		options = append(options, &option)
	}

	if err := rows.Err(); err != nil {
		r.log.WithError(err).Error("Error iterating option rows")
		return nil, fmt.Errorf("error iterating option rows: %w", err)
	}

	return options, nil
}
