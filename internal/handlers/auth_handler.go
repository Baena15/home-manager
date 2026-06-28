// Package handlers provides HTTP handlers for authentication.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gentleman-programming/home-manager/internal/config"
	"github.com/gentleman-programming/home-manager/internal/store"
	"github.com/gentleman-programming/home-manager/pkg/auth"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*store.User, error)
	Create(ctx context.Context, email, passwordHash, role string) (*store.User, error)
}

// ─── LoginRequest ───────────────────────────────────────────────────

// LoginRequest represents a user login attempt.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ─── LoginResponse ──────────────────────────────────────────────────

// LoginResponse represents a successful login response.
type LoginResponse struct {
	Token string     `json:"token"`
	User  store.User `json:"user"`
}

// ─── RegisterRequest ────────────────────────────────────────────────

// RegisterRequest represents a user registration attempt.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ─── RegisterResponse ───────────────────────────────────────────────

// RegisterResponse represents a successful registration response.
type RegisterResponse struct {
	Token string     `json:"token"`
	User  store.User `json:"user"`
}

// ─── AuthHandler ────────────────────────────────────────────────────

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	cfg   *config.Config
	users UserRepository
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cfg *config.Config, users UserRepository) *AuthHandler {
	return &AuthHandler{
		cfg:   cfg,
		users: users,
	}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "email and password are required")
		return
	}

	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "password must be at least 8 characters")
		return
	}

	_, err := h.users.GetByEmail(r.Context(), req.Email)
	if err == nil {
		respondError(w, http.StatusConflict, "CONFLICT", "email already registered")
		return
	}
	if !errors.Is(err, store.ErrUserNotFound) {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check email availability")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to process password")
		return
	}

	user, err := h.users.Create(r.Context(), req.Email, hash, "partner")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create user")
		return
	}

	token, err := auth.GenerateToken(h.cfg.JWTSecret, h.cfg.JWTExpiration, user.ID, user.Email, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate token")
		return
	}

	respondJSON(w, http.StatusCreated, RegisterResponse{
		Token: token,
		User:  *user,
	})
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "email and password are required")
		return
	}

	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "authentication failed")
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "invalid credentials")
		return
	}

	token, err := auth.GenerateToken(h.cfg.JWTSecret, h.cfg.JWTExpiration, user.ID, user.Email, user.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to generate token")
		return
	}

	respondJSON(w, http.StatusOK, LoginResponse{
		Token: token,
		User:  *user,
	})
}

// ─── Response helpers ───────────────────────────────────────────────

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error":     message,
		"code":      code,
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ErrorResponse is a standard API error response.
type ErrorResponse struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	Status    int    `json:"status"`
	Timestamp string `json:"timestamp"`
}
