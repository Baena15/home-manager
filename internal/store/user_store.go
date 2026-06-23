// Package store provides data access for users.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// ─── User ───────────────────────────────────────────────────────────

// User represents a household account.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// ─── UserStore ──────────────────────────────────────────────────────

// UserStore handles user data access.
type UserStore struct {
	db *DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *DB) *UserStore {
	return &UserStore{db: db}
}

// GetByEmail returns a user by email address.
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`

	row := s.db.QueryRowContext(ctx, query, email)

	user := &User{}
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Create inserts a new user and returns the created user.
func (s *UserStore) Create(ctx context.Context, email, passwordHash, role string) (*User, error) {
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role, created_at
	`

	user := &User{}
	if err := s.db.QueryRowContext(ctx, query, email, passwordHash, role).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Count returns the number of users in the database.
func (s *UserStore) Count(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
