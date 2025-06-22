// Package service provides session management business logic.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/octop162/normal-form-app-by-claude/internal/dto"
	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/internal/repository"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"

	"github.com/google/uuid"
)

const (
	// Default session timeout duration
	defaultSessionTimeout = 4 * time.Hour
)

// SessionService defines the interface for session business logic
type SessionService interface {
	CreateSession(ctx context.Context, req *dto.SessionCreateRequest) (*dto.SessionCreateResponse, error)
	GetSession(ctx context.Context, sessionID string) (*dto.SessionGetResponse, error)
	UpdateSession(ctx context.Context, sessionID string, req *dto.SessionUpdateRequest) (*dto.SessionUpdateResponse, error)
	DeleteSession(ctx context.Context, sessionID string) (*dto.SessionDeleteResponse, error)
	CleanupExpiredSessions(ctx context.Context) (int64, error)
	ExtendSession(ctx context.Context, sessionID string, duration time.Duration) (*dto.SessionUpdateResponse, error)
	IsSessionValid(ctx context.Context, sessionID string) (bool, error)
}

// sessionService implements SessionService
type sessionService struct {
	sessionRepo repository.SessionRepository
	log         *logger.Logger
}

// NewSessionService creates a new session service
func NewSessionService(
	sessionRepo repository.SessionRepository,
	log *logger.Logger,
) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		log:         log,
	}
}

// CreateSession creates a new session with user data
func (s *sessionService) CreateSession(
	ctx context.Context, req *dto.SessionCreateRequest,
) (*dto.SessionCreateResponse, error) {
	// Generate unique session ID
	sessionID := uuid.New().String()

	// Calculate expiration time
	expiresAt := time.Now().Add(defaultSessionTimeout)

	// Create session model
	session := &model.UserSession{
		ID:        sessionID,
		UserData:  req.UserData,
		ExpiresAt: expiresAt,
	}

	// Save session
	createdSession, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		s.log.WithError(err).Error("Failed to create session")
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.log.WithField("session_id", sessionID).Info("Session created successfully")

	return &dto.SessionCreateResponse{
		SessionID: createdSession.ID,
		ExpiresAt: createdSession.ExpiresAt,
	}, nil
}

// GetSession retrieves session data by ID
func (s *sessionService) GetSession(ctx context.Context, sessionID string) (*dto.SessionGetResponse, error) {
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		s.log.WithError(err).WithField("session_id", sessionID).Error("Failed to get session")
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		s.log.WithField("session_id", sessionID).Warn("Attempted to access expired session")
		return nil, fmt.Errorf("session has expired")
	}

	return &dto.SessionGetResponse{
		SessionID: session.ID,
		UserData:  session.UserData,
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

// UpdateSession updates session data and extends expiration
func (s *sessionService) UpdateSession(
	ctx context.Context, sessionID string, req *dto.SessionUpdateRequest,
) (*dto.SessionUpdateResponse, error) {
	// Get existing session
	existingSession, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if session is expired
	if existingSession.IsExpired() {
		return nil, fmt.Errorf("session has expired")
	}

	// Update session data and extend expiration
	existingSession.UserData = req.UserData
	existingSession.ExpiresAt = time.Now().Add(defaultSessionTimeout)

	// Save updated session
	updatedSession, err := s.sessionRepo.Update(ctx, existingSession)
	if err != nil {
		s.log.WithError(err).WithField("session_id", sessionID).Error("Failed to update session")
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	s.log.WithField("session_id", sessionID).Info("Session updated successfully")

	return &dto.SessionUpdateResponse{
		SessionID: updatedSession.ID,
		ExpiresAt: updatedSession.ExpiresAt,
		UpdatedAt: updatedSession.UpdatedAt,
	}, nil
}

// DeleteSession deletes a session
func (s *sessionService) DeleteSession(ctx context.Context, sessionID string) (*dto.SessionDeleteResponse, error) {
	err := s.sessionRepo.Delete(ctx, sessionID)
	if err != nil {
		s.log.WithError(err).WithField("session_id", sessionID).Error("Failed to delete session")
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	s.log.WithField("session_id", sessionID).Info("Session deleted successfully")

	return &dto.SessionDeleteResponse{
		Message: "Session deleted successfully",
	}, nil
}

// CleanupExpiredSessions removes all expired sessions
func (s *sessionService) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	deletedCount, err := s.sessionRepo.DeleteExpired(ctx)
	if err != nil {
		s.log.WithError(err).Error("Failed to cleanup expired sessions")
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	if deletedCount > 0 {
		s.log.WithField("deleted_count", deletedCount).Info("Expired sessions cleaned up")
	}

	return deletedCount, nil
}

// ExtendSession extends session expiration time
func (s *sessionService) ExtendSession(
	ctx context.Context, sessionID string, duration time.Duration,
) (*dto.SessionUpdateResponse, error) {
	// Get existing session
	existingSession, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Check if session is expired
	if existingSession.IsExpired() {
		return nil, fmt.Errorf("session has expired")
	}

	// Extend expiration time
	existingSession.ExpiresAt = time.Now().Add(duration)

	// Save updated session
	updatedSession, err := s.sessionRepo.Update(ctx, existingSession)
	if err != nil {
		s.log.WithError(err).WithField("session_id", sessionID).Error("Failed to extend session")
		return nil, fmt.Errorf("failed to extend session: %w", err)
	}

	s.log.WithField("session_id", sessionID).
		WithField("duration", duration).
		Info("Session extended successfully")

	return &dto.SessionUpdateResponse{
		SessionID: updatedSession.ID,
		ExpiresAt: updatedSession.ExpiresAt,
		UpdatedAt: updatedSession.UpdatedAt,
	}, nil
}

// IsSessionValid checks if a session exists and is not expired
func (s *sessionService) IsSessionValid(ctx context.Context, sessionID string) (bool, error) {
	exists, err := s.sessionRepo.Exists(ctx, sessionID)
	if err != nil {
		s.log.WithError(err).WithField("session_id", sessionID).Error("Failed to check session validity")
		return false, fmt.Errorf("failed to check session validity: %w", err)
	}

	return exists, nil
}
