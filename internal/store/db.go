// Package store provides data access for PostgreSQL.
package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// ─── DB ─────────────────────────────────────────────────────────────

// DB wraps sql.DB with project-specific helpers.
type DB struct {
	*sql.DB
}

// New opens a PostgreSQL connection from a connection string.
func New(connString string) (*DB, error) {
	sqlDB, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{sqlDB}, nil
}
