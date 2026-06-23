package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gentleman-programming/home-manager/pkg/auth"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"
	token, _ := auth.GenerateToken(secret, time.Hour, "user-1", "test@example.com", "owner")

	middleware := AuthMiddleware(secret)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := ClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("claims not found in context")
		}
		if claims.UserID != "user-1" {
			t.Errorf("UserID = %q, want %q", claims.UserID, "user-1")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"
	middleware := AuthMiddleware(secret)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	body := rec.Body.String()
	if body == "" {
		t.Fatal("expected error response body")
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secret := "this-is-a-very-long-secret-key-for-jwt-signing"
	token, _ := auth.GenerateToken(secret, -time.Hour, "user-1", "test@example.com", "owner")

	middleware := AuthMiddleware(secret)
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"valid bearer", "Bearer token123", "token123"},
		{"lowercase bearer", "bearer token123", "token123"},
		{"missing header", "", ""},
		{"wrong scheme", "Basic token123", ""},
		{"only one part", "Bearer", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			got := extractBearerToken(req)
			if got != tt.want {
				t.Errorf("extractBearerToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
