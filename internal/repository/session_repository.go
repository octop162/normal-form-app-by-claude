// Package repository provides session data access functionality.
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	Create(ctx context.Context, session *model.UserSession) (*model.UserSession, error)
	GetByID(ctx context.Context, id string) (*model.UserSession, error)
	Update(ctx context.Context, session *model.UserSession) (*model.UserSession, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) (int64, error)
	Exists(ctx context.Context, id string) (bool, error)
}

// sessionRepository implements SessionRepository
type sessionRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *sql.DB, log *logger.Logger) SessionRepository {
	return &sessionRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new session
func (r *sessionRepository) Create(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	userDataJSON, err := json.Marshal(session.UserData)
	if err != nil {
		r.log.WithError(err).Error("Failed to marshal user data")
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	query := `
		INSERT INTO user_sessions (id, user_data, expires_at)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at`

	var createdSession model.UserSession
	err = r.db.QueryRowContext(ctx, query, session.ID, userDataJSON, session.ExpiresAt).
		Scan(&createdSession.CreatedAt, &createdSession.UpdatedAt)

	if err != nil {
		r.log.WithError(err).WithField("session_id", session.ID).Error("Failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	createdSession.ID = session.ID
	createdSession.UserData = session.UserData
	createdSession.ExpiresAt = session.ExpiresAt

	r.log.WithField("session_id", createdSession.ID).Info("Session created successfully")
	return &createdSession, nil
}

// GetByID retrieves a session by ID
func (r *sessionRepository) GetByID(ctx context.Context, id string) (*model.UserSession, error) {
	query := `
		SELECT id, user_data, expires_at, created_at, updated_at
		FROM user_sessions
		WHERE id = $1 AND expires_at > NOW()`

	var session model.UserSession
	var userDataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &userDataJSON, &session.ExpiresAt,
		&session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired: %w", err)
		}
		r.log.WithError(err).WithField("session_id", id).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Unmarshal user data
	if err := json.Unmarshal(userDataJSON, &session.UserData); err != nil {
		r.log.WithError(err).WithField("session_id", id).Error("Failed to unmarshal user data")
		return nil, fmt.Errorf("failed to unmarshal user data: %w", err)
	}

	return &session, nil
}

// Update updates an existing session
func (r *sessionRepository) Update(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	userDataJSON, err := json.Marshal(session.UserData)
	if err != nil {
		r.log.WithError(err).Error("Failed to marshal user data")
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	query := `
		UPDATE user_sessions SET
			user_data = $2,
			expires_at = $3,
			updated_at = NOW()
		WHERE id = $1 AND expires_at > NOW()
		RETURNING updated_at`

	err = r.db.QueryRowContext(ctx, query, session.ID, userDataJSON, session.ExpiresAt).
		Scan(&session.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		r.log.WithError(err).WithField("session_id", session.ID).Error("Failed to update session")
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	r.log.WithField("session_id", session.ID).Info("Session updated successfully")
	return session, nil
}

// Delete deletes a session by ID
func (r *sessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.WithError(err).WithField("session_id", id).Error("Failed to delete session")
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	r.log.WithField("session_id", id).Info("Session deleted successfully")
	return nil
}

// DeleteExpired deletes all expired sessions
func (r *sessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM user_sessions WHERE expires_at <= NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.log.WithError(err).Error("Failed to delete expired sessions")
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		r.log.WithField("deleted_count", rowsAffected).Info("Expired sessions deleted successfully")
	}

	return rowsAffected, nil
}

// Exists checks if a session exists and is not expired
func (r *sessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_sessions WHERE id = $1 AND expires_at > NOW())`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		r.log.WithError(err).WithField("session_id", id).Error("Failed to check session existence")
		return false, fmt.Errorf("failed to check session existence: %w", err)
	}

	return exists, nil
}
