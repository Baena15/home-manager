// Package store provides data access for shopping lists.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── ShoppingList ───────────────────────────────────────────────────

// ShoppingList represents a shopping list.
type ShoppingList struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	CreatedBy      string    `json:"created_by"`
	EstimatedTotal float64   `json:"estimated_total"`
	ItemCount      int       `json:"item_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ─── ListStore ──────────────────────────────────────────────────────

// ListStore handles shopping list data access.
type ListStore struct {
	db *DB
}

// NewListStore creates a new ListStore.
func NewListStore(db *DB) *ListStore {
	return &ListStore{db: db}
}

// Create inserts a new shopping list.
func (s *ListStore) Create(ctx context.Context, name, createdBy string) (*ShoppingList, error) {
	query := `
		INSERT INTO shopping_lists (name, created_by)
		VALUES ($1, $2)
		RETURNING id, name, status, created_by, created_at, updated_at
	`

	list := &ShoppingList{}
	if err := s.db.QueryRowContext(ctx, query, name, createdBy).Scan(
		&list.ID, &list.Name, &list.Status, &list.CreatedBy, &list.CreatedAt, &list.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return list, nil
}

// GetByID returns a list by id with aggregated totals.
func (s *ListStore) GetByID(ctx context.Context, id string) (*ShoppingList, error) {
	query := `
		SELECT
			l.id, l.name, l.status, l.created_by,
			COALESCE(SUM(i.total), 0) AS estimated_total,
			COUNT(i.id) AS item_count,
			l.created_at, l.updated_at
		FROM shopping_lists l
		LEFT JOIN shopping_list_items i ON i.list_id = l.id
		WHERE l.id = $1
		GROUP BY l.id
	`

	list := &ShoppingList{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&list.ID, &list.Name, &list.Status, &list.CreatedBy,
		&list.EstimatedTotal, &list.ItemCount,
		&list.CreatedAt, &list.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	return list, nil
}

// List returns shopping lists ordered by most recent.
func (s *ListStore) List(ctx context.Context, limit, offset int) ([]ShoppingList, int, error) {
	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM shopping_lists").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count lists: %w", err)
	}

	query := `
		SELECT
			l.id, l.name, l.status, l.created_by,
			COALESCE(SUM(i.total), 0) AS estimated_total,
			COUNT(i.id) AS item_count,
			l.created_at, l.updated_at
		FROM shopping_lists l
		LEFT JOIN shopping_list_items i ON i.list_id = l.id
		GROUP BY l.id
		ORDER BY l.updated_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list lists: %w", err)
	}
	defer rows.Close()

	lists := []ShoppingList{}
	for rows.Next() {
		list := ShoppingList{}
		if err := rows.Scan(
			&list.ID, &list.Name, &list.Status, &list.CreatedBy,
			&list.EstimatedTotal, &list.ItemCount,
			&list.CreatedAt, &list.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan list: %w", err)
		}
		lists = append(lists, list)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate lists: %w", err)
	}

	return lists, total, nil
}

// Update modifies a list name.
func (s *ListStore) Update(ctx context.Context, id, name string) (*ShoppingList, error) {
	query := `
		UPDATE shopping_lists
		SET name = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, status, created_by, created_at, updated_at
	`

	list := &ShoppingList{}
	if err := s.db.QueryRowContext(ctx, query, id, name).Scan(
		&list.ID, &list.Name, &list.Status, &list.CreatedBy, &list.CreatedAt, &list.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return list, nil
}

// UpdateStatus changes a list status.
func (s *ListStore) UpdateStatus(ctx context.Context, id, status string) (*ShoppingList, error) {
	query := `
		UPDATE shopping_lists
		SET status = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, status, created_by, created_at, updated_at
	`

	list := &ShoppingList{}
	if err := s.db.QueryRowContext(ctx, query, id, status).Scan(
		&list.ID, &list.Name, &list.Status, &list.CreatedBy, &list.CreatedAt, &list.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to update list status: %w", err)
	}

	return list, nil
}

// Delete removes a list and its items.
func (s *ListStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM shopping_lists WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("list not found")
	}

	return nil
}
