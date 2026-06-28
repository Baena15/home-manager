// Package store provides data access for products.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ─── Product ────────────────────────────────────────────────────────

// Product represents a household product.
type Product struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Unit          string     `json:"unit"`
	Category      string     `json:"category,omitempty"`
	LatestPrice   *float64   `json:"latest_price,omitempty"`
	LatestStore   string     `json:"latest_store,omitempty"`
	LatestPriceAt *time.Time `json:"latest_price_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ─── ProductStore ───────────────────────────────────────────────────

// ProductStore handles product data access.
type ProductStore struct {
	db *DB
}

// NewProductStore creates a new ProductStore.
func NewProductStore(db *DB) *ProductStore {
	return &ProductStore{db: db}
}

// Create inserts a new product.
func (s *ProductStore) Create(ctx context.Context, name, unit, category string) (*Product, error) {
	query := `
		INSERT INTO products (name, unit, category)
		VALUES ($1, $2, $3)
		RETURNING id, name, unit, category, created_at, updated_at
	`

	product := &Product{}
	if err := s.db.QueryRowContext(ctx, query, name, unit, category).Scan(
		&product.ID, &product.Name, &product.Unit, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetByID returns a product by id including its latest price.
func (s *ProductStore) GetByID(ctx context.Context, id string) (*Product, error) {
	query := `
		SELECT
			p.id, p.name, p.unit, p.category,
			p.created_at, p.updated_at,
			latest.amount, latest.store, latest.recorded_at
		FROM products p
		LEFT JOIN LATERAL (
			SELECT amount, store, recorded_at
			FROM product_prices
			WHERE product_id = p.id
			ORDER BY recorded_at DESC
			LIMIT 1
		) latest ON true
		WHERE p.id = $1
	`

	product := &Product{}
	var latestPrice sql.NullFloat64
	var latestStore sql.NullString
	var latestAt sql.NullTime

	if err := s.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Unit, &product.Category,
		&product.CreatedAt, &product.UpdatedAt,
		&latestPrice, &latestStore, &latestAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if latestPrice.Valid {
		product.LatestPrice = &latestPrice.Float64
		product.LatestStore = latestStore.String
		product.LatestPriceAt = &latestAt.Time
	}

	return product, nil
}

// List returns paginated products with optional search by name.
func (s *ProductStore) List(ctx context.Context, search string, limit, offset int) ([]Product, int, error) {
	countQuery := `SELECT COUNT(*) FROM products WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')`
	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, search).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	query := `
		SELECT
			p.id, p.name, p.unit, p.category,
			p.created_at, p.updated_at,
			latest.amount, latest.store, latest.recorded_at
		FROM products p
		LEFT JOIN LATERAL (
			SELECT amount, store, recorded_at
			FROM product_prices
			WHERE product_id = p.id
			ORDER BY recorded_at DESC
			LIMIT 1
		) latest ON true
		WHERE ($1 = '' OR p.name ILIKE '%' || $1 || '%')
		ORDER BY p.name
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, search, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		product := Product{}
		var latestPrice sql.NullFloat64
		var latestStore sql.NullString
		var latestAt sql.NullTime

		if err := rows.Scan(
			&product.ID, &product.Name, &product.Unit, &product.Category,
			&product.CreatedAt, &product.UpdatedAt,
			&latestPrice, &latestStore, &latestAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}

		if latestPrice.Valid {
			product.LatestPrice = &latestPrice.Float64
			product.LatestStore = latestStore.String
			product.LatestPriceAt = &latestAt.Time
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate products: %w", err)
	}

	return products, total, nil
}

// Update modifies a product's name, unit, or category.
func (s *ProductStore) Update(ctx context.Context, id, name, unit, category string) (*Product, error) {
	query := `
		UPDATE products
		SET name = $2, unit = $3, category = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, unit, category, created_at, updated_at
	`

	product := &Product{}
	if err := s.db.QueryRowContext(ctx, query, id, name, unit, category).Scan(
		&product.ID, &product.Name, &product.Unit, &product.Category, &product.CreatedAt, &product.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

// Delete removes a product and its price history.
func (s *ProductStore) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// Exists checks if a product name already exists.
func (s *ProductStore) Exists(ctx context.Context, name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE name ILIKE $1)`
	if err := s.db.QueryRowContext(ctx, query, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}
