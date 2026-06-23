// Package middleware provides HTTP middleware for the API.
package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gentleman-programming/home-manager/pkg/auth"
)

// ─── Context keys ───────────────────────────────────────────────────

type contextKey string

const (
	// ClaimsKey is the context key for JWT claims.
	ClaimsKey contextKey = "claims"
)

// ─── AuthMiddleware ─────────────────────────────────────────────────

// AuthMiddleware validates JWT tokens on protected routes.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractBearerToken(r)
			if tokenString == "" {
				respondAuthError(w, "MISSING_TOKEN", "authorization token is required")
				return
			}

			claims, err := auth.ValidateToken(secret, tokenString)
			if err != nil {
				code := "INVALID_TOKEN"
				if strings.Contains(err.Error(), "token is expired") {
					code = "TOKEN_EXPIRED"
				}
				respondAuthError(w, code, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext returns JWT claims from the request context.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*auth.Claims)
	return claims, ok
}

// extractBearerToken extracts the token from the Authorization header.
func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func respondAuthError(w http.ResponseWriter, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	respondJSON(w, map[string]interface{}{
		"error":  message,
		"code":   code,
		"status": http.StatusUnauthorized,
	})
}

func respondJSON(w http.ResponseWriter, payload interface{}) {
	_ = json.NewEncoder(w).Encode(payload)
}
