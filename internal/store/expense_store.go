// Package store provides data access for expenses.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── Expense ────────────────────────────────────────────────────────

// Expense represents a household expense.
type Expense struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Visibility      string    `json:"visibility"`
	SplitPercentage float64   `json:"split_percentage"`
	ExpenseDate     time.Time `json:"expense_date"`
	IsRecurring     bool      `json:"is_recurring"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ─── ExpenseStore ───────────────────────────────────────────────────

// ExpenseStore handles expense data access.
type ExpenseStore struct {
	db *DB
}

// NewExpenseStore creates a new ExpenseStore.
func NewExpenseStore(db *DB) *ExpenseStore {
	return &ExpenseStore{db: db}
}

// Create inserts a new expense and returns it.
func (s *ExpenseStore) Create(ctx context.Context, userID string, amount float64, description, category, visibility string, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*Expense, error) {
	query := `
		INSERT INTO expenses (user_id, amount, description, category, visibility, split_percentage, expense_date, is_recurring)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, amount, description, category, visibility, split_percentage, expense_date, is_recurring, created_at, updated_at
	`

	expense := &Expense{}
	if err := s.db.QueryRowContext(ctx, query, userID, amount, description, category, visibility, splitPercentage, expenseDate, isRecurring).Scan(
		&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expense.Category,
		&expense.Visibility, &expense.SplitPercentage, &expense.ExpenseDate, &expense.IsRecurring,
		&expense.CreatedAt, &expense.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

// GetByID returns an expense by ID.
func (s *ExpenseStore) GetByID(ctx context.Context, id string) (*Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, visibility, split_percentage, expense_date, is_recurring, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`

	expense := &Expense{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expense.Category,
		&expense.Visibility, &expense.SplitPercentage, &expense.ExpenseDate, &expense.IsRecurring,
		&expense.CreatedAt, &expense.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expense not found")
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	return expense, nil
}

// List returns expenses visible to the given user, optionally filtered.
func (s *ExpenseStore) List(ctx context.Context, userID, visibility, from, to string) ([]Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, visibility, split_percentage, expense_date, is_recurring, created_at, updated_at
		FROM expenses
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
		query += fmt.Sprintf(" AND expense_date >= $%d", argCount)
		args = append(args, from)
	}
	if to != "" {
		argCount++
		query += fmt.Sprintf(" AND expense_date <= $%d", argCount)
		args = append(args, to)
	}

	query += " ORDER BY expense_date DESC, created_at DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list expenses: %w", err)
	}
	defer rows.Close()

	return scanExpenses(rows)
}

// Update modifies an existing expense.
func (s *ExpenseStore) Update(ctx context.Context, id, userID, description, category, visibility string, amount, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*Expense, error) {
	query := `
		UPDATE expenses
		SET description = $1, category = $2, visibility = $3, amount = $4, split_percentage = $5, expense_date = $6, is_recurring = $7, updated_at = NOW()
		WHERE id = $8 AND user_id = $9
		RETURNING id, user_id, amount, description, category, visibility, split_percentage, expense_date, is_recurring, created_at, updated_at
	`

	expense := &Expense{}
	if err := s.db.QueryRowContext(ctx, query, description, category, visibility, amount, splitPercentage, expenseDate, isRecurring, id, userID).Scan(
		&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expense.Category,
		&expense.Visibility, &expense.SplitPercentage, &expense.ExpenseDate, &expense.IsRecurring,
		&expense.CreatedAt, &expense.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("expense not found or not owned by user")
		}
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	return expense, nil
}

// Delete removes an expense owned by the given user.
func (s *ExpenseStore) Delete(ctx context.Context, id, userID string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM expenses WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("expense not found or not owned by user")
	}

	return nil
}

// scanExpenses scans rows into a slice of Expense.
func scanExpenses(rows *sql.Rows) ([]Expense, error) {
	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(
			&expense.ID, &expense.UserID, &expense.Amount, &expense.Description, &expense.Category,
			&expense.Visibility, &expense.SplitPercentage, &expense.ExpenseDate, &expense.IsRecurring,
			&expense.CreatedAt, &expense.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expenses: %w", err)
	}

	return expenses, nil
}
