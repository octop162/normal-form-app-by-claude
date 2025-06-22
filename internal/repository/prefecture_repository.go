// Package repository provides prefecture master data access functionality.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// PrefectureRepository defines the interface for prefecture master data access
type PrefectureRepository interface {
	GetAll(ctx context.Context) ([]*model.PrefectureMaster, error)
	GetByCode(ctx context.Context, prefectureCode string) (*model.PrefectureMaster, error)
	GetByName(ctx context.Context, prefectureName string) (*model.PrefectureMaster, error)
	GetByRegion(ctx context.Context, region string) ([]*model.PrefectureMaster, error)
	GetActive(ctx context.Context) ([]*model.PrefectureMaster, error)
}

// prefectureRepository implements PrefectureRepository
type prefectureRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewPrefectureRepository creates a new prefecture repository
func NewPrefectureRepository(db *sql.DB, log *logger.Logger) PrefectureRepository {
	return &prefectureRepository{
		db:  db,
		log: log,
	}
}

// GetAll retrieves all prefecture master data
func (r *prefectureRepository) GetAll(ctx context.Context) ([]*model.PrefectureMaster, error) {
	query := `
		SELECT id, prefecture_code, prefecture_name, region, is_active, created_at
		FROM prefectures_master
		ORDER BY prefecture_code ASC`

	return r.queryPrefectures(ctx, query)
}

// GetByCode retrieves a prefecture by prefecture code
func (r *prefectureRepository) GetByCode(ctx context.Context, prefectureCode string) (*model.PrefectureMaster, error) {
	query := `
		SELECT id, prefecture_code, prefecture_name, region, is_active, created_at
		FROM prefectures_master
		WHERE prefecture_code = $1`

	var prefecture model.PrefectureMaster
	err := r.db.QueryRowContext(ctx, query, prefectureCode).Scan(
		&prefecture.ID, &prefecture.PrefectureCode, &prefecture.PrefectureName,
		&prefecture.Region, &prefecture.IsActive, &prefecture.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("prefecture not found: %w", err)
		}
		r.log.WithError(err).WithField("prefecture_code", prefectureCode).Error("Failed to get prefecture by code")
		return nil, fmt.Errorf("failed to get prefecture by code: %w", err)
	}

	return &prefecture, nil
}

// GetByName retrieves a prefecture by prefecture name
func (r *prefectureRepository) GetByName(ctx context.Context, prefectureName string) (*model.PrefectureMaster, error) {
	query := `
		SELECT id, prefecture_code, prefecture_name, region, is_active, created_at
		FROM prefectures_master
		WHERE prefecture_name = $1`

	var prefecture model.PrefectureMaster
	err := r.db.QueryRowContext(ctx, query, prefectureName).Scan(
		&prefecture.ID, &prefecture.PrefectureCode, &prefecture.PrefectureName,
		&prefecture.Region, &prefecture.IsActive, &prefecture.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("prefecture not found: %w", err)
		}
		r.log.WithError(err).WithField("prefecture_name", prefectureName).Error("Failed to get prefecture by name")
		return nil, fmt.Errorf("failed to get prefecture by name: %w", err)
	}

	return &prefecture, nil
}

// GetByRegion retrieves prefectures by region
func (r *prefectureRepository) GetByRegion(ctx context.Context, region string) ([]*model.PrefectureMaster, error) {
	query := `
		SELECT id, prefecture_code, prefecture_name, region, is_active, created_at
		FROM prefectures_master
		WHERE region = $1
		ORDER BY prefecture_code ASC`

	rows, err := r.db.QueryContext(ctx, query, region)
	if err != nil {
		r.log.WithError(err).WithField("region", region).Error("Failed to get prefectures by region")
		return nil, fmt.Errorf("failed to get prefectures by region: %w", err)
	}
	defer rows.Close()

	return r.scanPrefectures(rows)
}

// GetActive retrieves all active prefectures
func (r *prefectureRepository) GetActive(ctx context.Context) ([]*model.PrefectureMaster, error) {
	query := `
		SELECT id, prefecture_code, prefecture_name, region, is_active, created_at
		FROM prefectures_master
		WHERE is_active = true
		ORDER BY prefecture_code ASC`

	return r.queryPrefectures(ctx, query)
}

// queryPrefectures executes a query and returns prefectures
func (r *prefectureRepository) queryPrefectures(
	ctx context.Context, query string, args ...any,
) ([]*model.PrefectureMaster, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.log.WithError(err).Error("Failed to query prefectures")
		return nil, fmt.Errorf("failed to query prefectures: %w", err)
	}
	defer rows.Close()

	return r.scanPrefectures(rows)
}

// scanPrefectures scans rows into prefecture master structs
func (r *prefectureRepository) scanPrefectures(rows *sql.Rows) ([]*model.PrefectureMaster, error) {
	var prefectures []*model.PrefectureMaster

	for rows.Next() {
		var prefecture model.PrefectureMaster
		err := rows.Scan(
			&prefecture.ID, &prefecture.PrefectureCode, &prefecture.PrefectureName,
			&prefecture.Region, &prefecture.IsActive, &prefecture.CreatedAt,
		)
		if err != nil {
			r.log.WithError(err).Error("Failed to scan prefecture row")
			return nil, fmt.Errorf("failed to scan prefecture row: %w", err)
		}
		prefectures = append(prefectures, &prefecture)
	}

	if err := rows.Err(); err != nil {
		r.log.WithError(err).Error("Error iterating prefecture rows")
		return nil, fmt.Errorf("error iterating prefecture rows: %w", err)
	}

	return prefectures, nil
}
