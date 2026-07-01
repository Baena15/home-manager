// Package store provides data access for settlements.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── Settlement ─────────────────────────────────────────────────────

// Settlement represents a payment between two household users.
type Settlement struct {
	ID             string    `json:"id"`
	FromUserID     string    `json:"from_user_id"`
	FromUserEmail  string    `json:"from_user_email"`
	ToUserID       string    `json:"to_user_id"`
	ToUserEmail    string    `json:"to_user_email"`
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
	SettlementDate time.Time `json:"settlement_date"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ─── SettlementStore ────────────────────────────────────────────────

// SettlementStore handles settlement data access.
type SettlementStore struct {
	db *DB
}

// NewSettlementStore creates a new SettlementStore.
func NewSettlementStore(db *DB) *SettlementStore {
	return &SettlementStore{db: db}
}

// Create inserts a new settlement and returns it.
func (s *SettlementStore) Create(ctx context.Context, fromUserID, toUserID string, amount float64, description string, settlementDate time.Time) (*Settlement, error) {
	query := `
		INSERT INTO settlements (from_user_id, to_user_id, amount, description, settlement_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, from_user_id, to_user_id, amount, description, settlement_date, created_at, updated_at
	`

	settlement := &Settlement{}
	if err := s.db.QueryRowContext(ctx, query, fromUserID, toUserID, amount, description, settlementDate).Scan(
		&settlement.ID, &settlement.FromUserID, &settlement.ToUserID, &settlement.Amount,
		&settlement.Description, &settlement.SettlementDate, &settlement.CreatedAt, &settlement.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create settlement: %w", err)
	}

	return settlement, nil
}

// List returns settlements where the given user is either payer or receiver.
func (s *SettlementStore) List(ctx context.Context, userID, from, to string) ([]Settlement, error) {
	query := `
		SELECT s.id, s.from_user_id, fu.email, s.to_user_id, tu.email, s.amount, s.description, s.settlement_date, s.created_at, s.updated_at
		FROM settlements s
		JOIN users fu ON fu.id = s.from_user_id
		JOIN users tu ON tu.id = s.to_user_id
		WHERE (s.from_user_id = $1 OR s.to_user_id = $1)
	`
	args := []interface{}{userID}
	argCount := 1

	if from != "" {
		argCount++
		query += fmt.Sprintf(" AND settlement_date >= $%d", argCount)
		args = append(args, from)
	}
	if to != "" {
		argCount++
		query += fmt.Sprintf(" AND settlement_date <= $%d", argCount)
		args = append(args, to)
	}

	query += " ORDER BY settlement_date DESC, created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list settlements: %w", err)
	}
	defer rows.Close()

	return scanSettlements(rows)
}

// GetByID returns a settlement by ID.
func (s *SettlementStore) GetByID(ctx context.Context, id string) (*Settlement, error) {
	query := `
		SELECT s.id, s.from_user_id, fu.email, s.to_user_id, tu.email, s.amount, s.description, s.settlement_date, s.created_at, s.updated_at
		FROM settlements s
		JOIN users fu ON fu.id = s.from_user_id
		JOIN users tu ON tu.id = s.to_user_id
		WHERE s.id = $1
	`

	settlement := &Settlement{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&settlement.ID, &settlement.FromUserID, &settlement.FromUserEmail, &settlement.ToUserID, &settlement.ToUserEmail, &settlement.Amount,
		&settlement.Description, &settlement.SettlementDate, &settlement.CreatedAt, &settlement.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("settlement not found")
		}
		return nil, fmt.Errorf("failed to get settlement: %w", err)
	}

	return settlement, nil
}

// Delete removes a settlement.
func (s *SettlementStore) Delete(ctx context.Context, id, userID string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM settlements WHERE id = $1 AND (from_user_id = $2 OR to_user_id = $2)", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete settlement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("settlement not found or not involved user")
	}

	return nil
}

// scanSettlements scans rows into a slice of Settlement.
func scanSettlements(rows *sql.Rows) ([]Settlement, error) {
	var settlements []Settlement
	for rows.Next() {
		var settlement Settlement
		if err := rows.Scan(
			&settlement.ID, &settlement.FromUserID, &settlement.FromUserEmail, &settlement.ToUserID, &settlement.ToUserEmail, &settlement.Amount,
			&settlement.Description, &settlement.SettlementDate, &settlement.CreatedAt, &settlement.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan settlement: %w", err)
		}
		settlements = append(settlements, settlement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settlements: %w", err)
	}

	return settlements, nil
}
