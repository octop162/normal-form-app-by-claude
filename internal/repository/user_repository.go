// Package repository provides data access layer interfaces and implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/octop162/normal-form-app-by-claude/internal/model"
	"github.com/octop162/normal-form-app-by-claude/pkg/logger"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, error)
}

// userRepository implements UserRepository
type userRepository struct {
	db  *sql.DB
	log *logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB, log *logger.Logger) UserRepository {
	return &userRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (
			last_name, first_name, last_name_kana, first_name_kana,
			phone1, phone2, phone3, postal_code1, postal_code2,
			prefecture, city, town, chome, banchi, go, building, room,
			email, plan_type
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		) RETURNING id, created_at, updated_at`

	var createdUser model.User
	err := r.db.QueryRowContext(ctx, query,
		user.LastName, user.FirstName, user.LastNameKana, user.FirstNameKana,
		user.Phone1, user.Phone2, user.Phone3, user.PostalCode1, user.PostalCode2,
		user.Prefecture, user.City, user.Town, user.Chome, user.Banchi,
		user.Go, user.Building, user.Room, user.Email, user.PlanType,
	).Scan(&createdUser.ID, &createdUser.CreatedAt, &createdUser.UpdatedAt)

	if err != nil {
		r.log.WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Copy all fields from input user to created user
	createdUser.LastName = user.LastName
	createdUser.FirstName = user.FirstName
	createdUser.LastNameKana = user.LastNameKana
	createdUser.FirstNameKana = user.FirstNameKana
	createdUser.Phone1 = user.Phone1
	createdUser.Phone2 = user.Phone2
	createdUser.Phone3 = user.Phone3
	createdUser.PostalCode1 = user.PostalCode1
	createdUser.PostalCode2 = user.PostalCode2
	createdUser.Prefecture = user.Prefecture
	createdUser.City = user.City
	createdUser.Town = user.Town
	createdUser.Chome = user.Chome
	createdUser.Banchi = user.Banchi
	createdUser.Go = user.Go
	createdUser.Building = user.Building
	createdUser.Room = user.Room
	createdUser.Email = user.Email
	createdUser.PlanType = user.PlanType

	r.log.WithField("user_id", createdUser.ID).Info("User created successfully")
	return &createdUser, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT id, last_name, first_name, last_name_kana, first_name_kana,
			   phone1, phone2, phone3, postal_code1, postal_code2,
			   prefecture, city, town, chome, banchi, go, building, room,
			   email, plan_type, created_at, updated_at
		FROM users WHERE id = $1`

	user, err := r.scanSingleUser(ctx, query, id)
	if err != nil {
		r.log.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, last_name, first_name, last_name_kana, first_name_kana,
			   phone1, phone2, phone3, postal_code1, postal_code2,
			   prefecture, city, town, chome, banchi, go, building, room,
			   email, plan_type, created_at, updated_at
		FROM users WHERE email = $1`

	user, err := r.scanSingleUser(ctx, query, email)
	if err != nil {
		r.log.WithError(err).WithField("email", email).Error("Failed to get user by email")
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// scanSingleUser scans a single user from query result
func (r *userRepository) scanSingleUser(ctx context.Context, query string, arg any) (*model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID, &user.LastName, &user.FirstName, &user.LastNameKana, &user.FirstNameKana,
		&user.Phone1, &user.Phone2, &user.Phone3, &user.PostalCode1, &user.PostalCode2,
		&user.Prefecture, &user.City, &user.Town, &user.Chome, &user.Banchi,
		&user.Go, &user.Building, &user.Room, &user.Email, &user.PlanType,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		UPDATE users SET
			last_name = $2, first_name = $3, last_name_kana = $4, first_name_kana = $5,
			phone1 = $6, phone2 = $7, phone3 = $8, postal_code1 = $9, postal_code2 = $10,
			prefecture = $11, city = $12, town = $13, chome = $14, banchi = $15,
			go = $16, building = $17, room = $18, email = $19, plan_type = $20,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.LastName, user.FirstName, user.LastNameKana, user.FirstNameKana,
		user.Phone1, user.Phone2, user.Phone3, user.PostalCode1, user.PostalCode2,
		user.Prefecture, user.City, user.Town, user.Chome, user.Banchi,
		user.Go, user.Building, user.Room, user.Email, user.PlanType,
	).Scan(&user.UpdatedAt)

	if err != nil {
		r.log.WithError(err).WithField("user_id", user.ID).Error("Failed to update user")
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	r.log.WithField("user_id", user.ID).Info("User updated successfully")
	return user, nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.log.WithError(err).WithField("user_id", id).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.log.WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// ExistsByEmail checks if a user exists by email
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		r.log.WithError(err).WithField("email", email).Error("Failed to check user existence")
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// List retrieves a list of users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, last_name, first_name, last_name_kana, first_name_kana,
			   phone1, phone2, phone3, postal_code1, postal_code2,
			   prefecture, city, town, chome, banchi, go, building, room,
			   email, plan_type, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.log.WithError(err).Error("Failed to list users")
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		scanErr := rows.Scan(
			&user.ID, &user.LastName, &user.FirstName, &user.LastNameKana, &user.FirstNameKana,
			&user.Phone1, &user.Phone2, &user.Phone3, &user.PostalCode1, &user.PostalCode2,
			&user.Prefecture, &user.City, &user.Town, &user.Chome, &user.Banchi,
			&user.Go, &user.Building, &user.Room, &user.Email, &user.PlanType,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if scanErr != nil {
			r.log.WithError(scanErr).Error("Failed to scan user row")
			return nil, fmt.Errorf("failed to scan user row: %w", scanErr)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		r.log.WithError(err).Error("Error iterating user rows")
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
