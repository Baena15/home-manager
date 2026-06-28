// Package store provides database migration utilities.
package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all pending .sql migration files found in the migrations directory.
// Applied migrations are tracked in the schema_migrations table.
func (db *DB) RunMigrations(ctx context.Context) error {
	if err := db.ensureMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	files, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sqlFiles := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, name := range sqlFiles {
		applied, err := db.isMigrationApplied(ctx, name)
		if err != nil {
			return fmt.Errorf("failed to check migration status %s: %w", name, err)
		}
		if applied {
			continue
		}

		content, err := migrationsFS.ReadFile(filepath.Join("migrations", name))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		if err := db.runMigration(ctx, name, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", name, err)
		}
	}

	return nil
}

// ensureMigrationsTable creates the migration tracking table if it does not exist.
func (db *DB) ensureMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

// isMigrationApplied returns true if the given migration filename has already been applied.
func (db *DB) isMigrationApplied(ctx context.Context, filename string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE filename = $1)`
	if err := db.QueryRowContext(ctx, query, filename).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to query migration %s: %w", filename, err)
	}
	return exists, nil
}

// runMigration executes a single migration inside a transaction and records it.
func (db *DB) runMigration(ctx context.Context, filename, content string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.ExecContext(ctx, content); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to execute migration sql: %w", err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", filename); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}
