// Package store provides data access for shopping list items.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── ShoppingListItem ───────────────────────────────────────────────

// ShoppingListItem represents an item inside a shopping list.
type ShoppingListItem struct {
	ID          string    `json:"id"`
	ListID      string    `json:"list_id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name,omitempty"`
	Quantity    float64   `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Total       float64   `json:"total"`
	Purchased   bool      `json:"purchased"`
	CreatedAt   time.Time `json:"created_at"`
}

// ─── ListItemStore ──────────────────────────────────────────────────

// ListItemStore handles shopping list item data access.
type ListItemStore struct {
	db *DB
}

// NewListItemStore creates a new ListItemStore.
func NewListItemStore(db *DB) *ListItemStore {
	return &ListItemStore{db: db}
}

// Create inserts a new item into a list.
func (s *ListItemStore) Create(ctx context.Context, listID, productID string, quantity, unitPrice float64) (*ShoppingListItem, error) {
	total := quantity * unitPrice

	query := `
		INSERT INTO shopping_list_items (list_id, product_id, quantity, unit_price, total)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, list_id, product_id, quantity, unit_price, total, purchased, created_at
	`

	item := &ShoppingListItem{}
	if err := s.db.QueryRowContext(ctx, query, listID, productID, quantity, unitPrice, total).Scan(
		&item.ID, &item.ListID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Total, &item.Purchased, &item.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create list item: %w", err)
	}

	return item, nil
}

// ListByList returns items for a shopping list.
func (s *ListItemStore) ListByList(ctx context.Context, listID string) ([]ShoppingListItem, error) {
	query := `
		SELECT
			i.id, i.list_id, i.product_id, p.name,
			i.quantity, i.unit_price, i.total, i.purchased, i.created_at
		FROM shopping_list_items i
		JOIN products p ON p.id = i.product_id
		WHERE i.list_id = $1
		ORDER BY i.created_at
	`

	rows, err := s.db.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	items := []ShoppingListItem{}
	for rows.Next() {
		item := ShoppingListItem{}
		if err := rows.Scan(
			&item.ID, &item.ListID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPrice, &item.Total, &item.Purchased, &item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate items: %w", err)
	}

	return items, nil
}

// UpdatePurchased marks an item as purchased or not.
func (s *ListItemStore) UpdatePurchased(ctx context.Context, id string, purchased bool) (*ShoppingListItem, error) {
	query := `
		UPDATE shopping_list_items
		SET purchased = $2
		WHERE id = $1
		RETURNING id, list_id, product_id, quantity, unit_price, total, purchased, created_at
	`

	item := &ShoppingListItem{}
	if err := s.db.QueryRowContext(ctx, query, id, purchased).Scan(
		&item.ID, &item.ListID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Total, &item.Purchased, &item.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return item, nil
}

// Delete removes an item from a list.
func (s *ListItemStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM shopping_list_items WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// GetByID returns an item by id.
func (s *ListItemStore) GetByID(ctx context.Context, id string) (*ShoppingListItem, error) {
	query := `
		SELECT
			i.id, i.list_id, i.product_id, p.name,
			i.quantity, i.unit_price, i.total, i.purchased, i.created_at
		FROM shopping_list_items i
		JOIN products p ON p.id = i.product_id
		WHERE i.id = $1
	`

	item := &ShoppingListItem{}
	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.ListID, &item.ProductID, &item.ProductName,
		&item.Quantity, &item.UnitPrice, &item.Total, &item.Purchased, &item.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}
