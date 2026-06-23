// Package store provides database seeding utilities.
package store

import (
	"context"
	"fmt"

	"github.com/gentleman-programming/home-manager/pkg/auth"
)

// ─── SeedUsers ──────────────────────────────────────────────────────

// SeedUsers creates the two household users if the users table is empty.
func SeedUsers(ctx context.Context, users *UserStore, ownerEmail, ownerPassword, partnerEmail, partnerPassword string) error {
	count, err := users.Count(ctx)
	if err != nil {
		return fmt.Errorf("failed to check user count: %w", err)
	}

	if count > 0 {
		return nil
	}

	credentials := []struct {
		email    string
		password string
		role     string
	}{
		{email: ownerEmail, password: ownerPassword, role: "owner"},
		{email: partnerEmail, password: partnerPassword, role: "partner"},
	}

	for _, c := range credentials {
		hash, err := auth.HashPassword(c.password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", c.email, err)
		}

		if _, err := users.Create(ctx, c.email, hash, c.role); err != nil {
			return fmt.Errorf("failed to create user %s: %w", c.email, err)
		}
	}

	return nil
}
