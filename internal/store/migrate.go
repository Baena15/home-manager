// Package store provides database migration utilities.
package store

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all .sql migration files found in the migrations directory.
func (db *DB) RunMigrations(ctx context.Context) error {
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
		content, err := migrationsFS.ReadFile(filepath.Join("migrations", name))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", name, err)
		}
	}

	return nil
}
