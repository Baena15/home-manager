// Package store provides data access for product prices.
package store

import (
	"context"
	"fmt"
	"time"
)

// ─── ProductPrice ───────────────────────────────────────────────────

// ProductPrice represents a price record for a product.
type ProductPrice struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	Store      string    `json:"store"`
	Amount     float64   `json:"amount"`
	RecordedAt time.Time `json:"recorded_at"`
}

// ─── PriceStore ─────────────────────────────────────────────────────

// PriceStore handles product price data access.
type PriceStore struct {
	db *DB
}

// NewPriceStore creates a new PriceStore.
func NewPriceStore(db *DB) *PriceStore {
	return &PriceStore{db: db}
}

// Create inserts a new price record for a product.
func (s *PriceStore) Create(ctx context.Context, productID, store string, amount float64) (*ProductPrice, error) {
	query := `
		INSERT INTO product_prices (product_id, store, amount)
		VALUES ($1, $2, $3)
		RETURNING id, product_id, store, amount, recorded_at
	`

	price := &ProductPrice{}
	if err := s.db.QueryRowContext(ctx, query, productID, store, amount).Scan(
		&price.ID, &price.ProductID, &price.Store, &price.Amount, &price.RecordedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create price: %w", err)
	}

	return price, nil
}

// ListByProduct returns price history for a product ordered by most recent.
func (s *PriceStore) ListByProduct(ctx context.Context, productID string, limit int) ([]ProductPrice, error) {
	query := `
		SELECT id, product_id, store, amount, recorded_at
		FROM product_prices
		WHERE product_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := s.db.QueryContext(ctx, query, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list prices: %w", err)
	}
	defer rows.Close()

	prices := []ProductPrice{}
	for rows.Next() {
		price := ProductPrice{}
		if err := rows.Scan(&price.ID, &price.ProductID, &price.Store, &price.Amount, &price.RecordedAt); err != nil {
			return nil, fmt.Errorf("failed to scan price: %w", err)
		}
		prices = append(prices, price)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate prices: %w", err)
	}

	return prices, nil
}
