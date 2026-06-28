package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gentleman-programming/home-manager/internal/config"
	"github.com/gentleman-programming/home-manager/internal/store"
	"github.com/gentleman-programming/home-manager/pkg/auth"
)

type mockUserRepository struct {
	user *store.User
	err  error
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*store.User, error) {
	return m.user, m.err
}

func (m *mockUserRepository) Create(ctx context.Context, email, passwordHash, role string) (*store.User, error) {
	return m.user, m.err
}

func TestAuthHandler_Login_Success(t *testing.T) {
	password := "password123"
	hash, _ := auth.HashPassword(password)

	repo := &mockUserRepository{
		user: &store.User{
			ID:           "user-1",
			Email:        "owner@home.local",
			PasswordHash: hash,
			Role:         "owner",
			CreatedAt:    time.Now(),
		},
	}

	cfg := &config.Config{
		JWTSecret:     "this-is-a-very-long-secret-key-for-jwt-signing",
		JWTExpiration: time.Hour,
	}

	handler := NewAuthHandler(cfg, repo)

	body, _ := json.Marshal(LoginRequest{Email: "owner@home.local", Password: password})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Token == "" {
		t.Fatal("token is empty")
	}

	if resp.User.Email != "owner@home.local" {
		t.Errorf("email = %q, want %q", resp.User.Email, "owner@home.local")
	}
}

func TestAuthHandler_Login_InvalidPassword(t *testing.T) {
	password := "password123"
	hash, _ := auth.HashPassword(password)

	repo := &mockUserRepository{
		user: &store.User{
			ID:           "user-1",
			Email:        "owner@home.local",
			PasswordHash: hash,
			Role:         "owner",
			CreatedAt:    time.Now(),
		},
	}

	cfg := &config.Config{
		JWTSecret:     "this-is-a-very-long-secret-key-for-jwt-signing",
		JWTExpiration: time.Hour,
	}

	handler := NewAuthHandler(cfg, repo)

	body, _ := json.Marshal(LoginRequest{Email: "owner@home.local", Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepository{
		err: store.ErrUserNotFound,
	}

	cfg := &config.Config{
		JWTSecret:     "this-is-a-very-long-secret-key-for-jwt-signing",
		JWTExpiration: time.Hour,
	}

	handler := NewAuthHandler(cfg, repo)

	body, _ := json.Marshal(LoginRequest{Email: "unknown@home.local", Password: "password"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	repo := &mockUserRepository{}
	cfg := &config.Config{
		JWTSecret:     "this-is-a-very-long-secret-key-for-jwt-signing",
		JWTExpiration: time.Hour,
	}
	handler := NewAuthHandler(cfg, repo)

	body, _ := json.Marshal(LoginRequest{Email: "", Password: ""})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
