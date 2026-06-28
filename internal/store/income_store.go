// Package store provides data access for incomes.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── Income ─────────────────────────────────────────────────────────

// Income represents a household income.
type Income struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Visibility  string    `json:"visibility"`
	IncomeDate  time.Time `json:"income_date"`
	IsRecurring bool      `json:"is_recurring"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ─── IncomeStore ────────────────────────────────────────────────────

// IncomeStore handles income data access.
type IncomeStore struct {
	db *DB
}

// NewIncomeStore creates a new IncomeStore.
func NewIncomeStore(db *DB) *IncomeStore {
	return &IncomeStore{db: db}
}

// Create inserts a new income and returns it.
func (s *IncomeStore) Create(ctx context.Context, userID string, amount float64, description, category, visibility string, incomeDate time.Time, isRecurring bool) (*Income, error) {
	query := `
		INSERT INTO incomes (user_id, amount, description, category, visibility, income_date, is_recurring)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, amount, description, category, visibility, income_date, is_recurring, created_at, updated_at
	`

	income := &Income{}
	if err := s.db.QueryRowContext(ctx, query, userID, amount, description, category, visibility, incomeDate, isRecurring).Scan(
		&income.ID, &income.UserID, &income.Amount, &income.Description, &income.Category,
		&income.Visibility, &income.IncomeDate, &income.IsRecurring, &income.CreatedAt, &income.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create income: %w", err)
	}

	return income, nil
}

// GetByID returns an income by ID.
func (s *IncomeStore) GetByID(ctx context.Context, id string) (*Income, error) {
	query := `
		SELECT id, user_id, amount, description, category, visibility, income_date, is_recurring, created_at, updated_at
		FROM incomes
		WHERE id = $1
	`

	income := &Income{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&income.ID, &income.UserID, &income.Amount, &income.Description, &income.Category,
		&income.Visibility, &income.IncomeDate, &income.IsRecurring, &income.CreatedAt, &income.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("income not found")
		}
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	return income, nil
}

// List returns incomes visible to the given user, optionally filtered.
func (s *IncomeStore) List(ctx context.Context, userID, visibility, from, to string) ([]Income, error) {
	query := `
		SELECT id, user_id, amount, description, category, visibility, income_date, is_recurring, created_at, updated_at
		FROM incomes
		WHERE (user_id = $1 OR visibility = 'shared')
	`
	args := []interface{}{userID}
	argCount := 1

	if visibility != "" {
		argCount++
		query += fmt.Sprintf(" AND visibility = $%d", argCount)
		args = append(args, visibility)
	}
	if from != "" {
		argCount++
		query += fmt.Sprintf(" AND income_date >= $%d", argCount)
		args = append(args, from)
	}
	if to != "" {
		argCount++
		query += fmt.Sprintf(" AND income_date <= $%d", argCount)
		args = append(args, to)
	}

	query += " ORDER BY income_date DESC, created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list incomes: %w", err)
	}
	defer rows.Close()

	return scanIncomes(rows)
}

// Update modifies an existing income.
func (s *IncomeStore) Update(ctx context.Context, id, userID, description, category, visibility string, amount float64, incomeDate time.Time, isRecurring bool) (*Income, error) {
	query := `
		UPDATE incomes
		SET description = $1, category = $2, visibility = $3, amount = $4, income_date = $5, is_recurring = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING id, user_id, amount, description, category, visibility, income_date, is_recurring, created_at, updated_at
	`

	income := &Income{}
	if err := s.db.QueryRowContext(ctx, query, description, category, visibility, amount, incomeDate, isRecurring, id, userID).Scan(
		&income.ID, &income.UserID, &income.Amount, &income.Description, &income.Category,
		&income.Visibility, &income.IncomeDate, &income.IsRecurring, &income.CreatedAt, &income.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("income not found or not owned by user")
		}
		return nil, fmt.Errorf("failed to update income: %w", err)
	}

	return income, nil
}

// Delete removes an income owned by the given user.
func (s *IncomeStore) Delete(ctx context.Context, id, userID string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM incomes WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete income: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("income not found or not owned by user")
	}

	return nil
}

// scanIncomes scans rows into a slice of Income.
func scanIncomes(rows *sql.Rows) ([]Income, error) {
	var incomes []Income
	for rows.Next() {
		var income Income
		if err := rows.Scan(
			&income.ID, &income.UserID, &income.Amount, &income.Description, &income.Category,
			&income.Visibility, &income.IncomeDate, &income.IsRecurring, &income.CreatedAt, &income.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan income: %w", err)
		}
		incomes = append(incomes, income)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating incomes: %w", err)
	}

	return incomes, nil
}
