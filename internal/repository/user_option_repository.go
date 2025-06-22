// Package repository provides user option data access functionality.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// UserOptionRepository defines the interface for user option data access
type UserOptionRepository interface {
	Create(ctx context.Context, userOption *model.UserOption) (*model.UserOption, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.UserOption, error)
	DeleteByUserID(ctx context.Context, userID int) error
	CreateBatch(ctx context.Context, userOptions []*model.UserOption) error
	DeleteByUserIDAndOptionType(ctx context.Context, userID int, optionType string) error
}

// userOptionRepository implements UserOptionRepository
type userOptionRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewUserOptionRepository creates a new user option repository
func NewUserOptionRepository(db *sql.DB, log *logger.Logger) UserOptionRepository {
	return &userOptionRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new user option
func (r *userOptionRepository) Create(ctx context.Context, userOption *model.UserOption) (*model.UserOption, error) {
	query := `
		INSERT INTO user_options (user_id, option_type)
		VALUES ($1, $2)
		RETURNING id, created_at`

	var createdOption model.UserOption
	err := r.db.QueryRowContext(ctx, query, userOption.UserID, userOption.OptionType).
		Scan(&createdOption.ID, &createdOption.CreatedAt)

	if err != nil {
		r.log.WithError(err).
			WithField("user_id", userOption.UserID).
			WithField("option_type", userOption.OptionType).
			Error("Failed to create user option")
		return nil, fmt.Errorf("failed to create user option: %w", err)
	}

	createdOption.UserID = userOption.UserID
	createdOption.OptionType = userOption.OptionType

	r.log.WithField("user_option_id", createdOption.ID).Info("User option created successfully")
	return &createdOption, nil
}

// GetByUserID retrieves all user options by user ID
func (r *userOptionRepository) GetByUserID(ctx context.Context, userID int) ([]*model.UserOption, error) {
	query := `
		SELECT id, user_id, option_type, created_at
		FROM user_options
		WHERE user_id = $1
		ORDER BY created_at ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("Failed to get user options")
		return nil, fmt.Errorf("failed to get user options: %w", err)
	}
	defer rows.Close()

	var userOptions []*model.UserOption
	for rows.Next() {
		var option model.UserOption
		scanErr := rows.Scan(&option.ID, &option.UserID, &option.OptionType, &option.CreatedAt)
		if scanErr != nil {
			r.log.WithError(scanErr).Error("Failed to scan user option row")
			return nil, fmt.Errorf("failed to scan user option row: %w", scanErr)
		}
		userOptions = append(userOptions, &option)
	}

	if err = rows.Err(); err != nil {
		r.log.WithError(err).Error("Error iterating user option rows")
		return nil, fmt.Errorf("error iterating user option rows: %w", err)
	}

	return userOptions, nil
}

// DeleteByUserID deletes all user options by user ID
func (r *userOptionRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM user_options WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.WithError(err).WithField("user_id", userID).Error("Failed to delete user options")
		return fmt.Errorf("failed to delete user options: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.log.WithField("user_id", userID).
		WithField("deleted_count", rowsAffected).
		Info("User options deleted successfully")
	return nil
}

// CreateBatch creates multiple user options in a single transaction
func (r *userOptionRepository) CreateBatch(ctx context.Context, userOptions []*model.UserOption) error {
	if len(userOptions) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.log.WithError(rollbackErr).Error("Failed to rollback transaction")
			}
		}
	}()

	query := `INSERT INTO user_options (user_id, option_type) VALUES ($1, $2)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, option := range userOptions {
		_, err = stmt.ExecContext(ctx, option.UserID, option.OptionType)
		if err != nil {
			r.log.WithError(err).
				WithField("user_id", option.UserID).
				WithField("option_type", option.OptionType).
				Error("Failed to insert user option in batch")
			return fmt.Errorf("failed to insert user option: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	r.log.WithField("batch_size", len(userOptions)).Info("User options batch created successfully")
	return nil
}

// DeleteByUserIDAndOptionType deletes a specific user option
func (r *userOptionRepository) DeleteByUserIDAndOptionType(ctx context.Context, userID int, optionType string) error {
	query := `DELETE FROM user_options WHERE user_id = $1 AND option_type = $2`

	result, err := r.db.ExecContext(ctx, query, userID, optionType)
	if err != nil {
		r.log.WithError(err).
			WithField("user_id", userID).
			WithField("option_type", optionType).
			Error("Failed to delete user option")
		return fmt.Errorf("failed to delete user option: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user option not found")
	}

	r.log.WithField("user_id", userID).
		WithField("option_type", optionType).
		Info("User option deleted successfully")
	return nil
}
