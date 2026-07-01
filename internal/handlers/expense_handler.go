// Package handlers provides HTTP handlers for expenses.
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── ExpenseRepository ──────────────────────────────────────────────

// ExpenseRepository defines the interface for expense data access.
type ExpenseRepository interface {
	Create(ctx context.Context, userID string, amount float64, description, category, visibility string, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*store.Expense, error)
	GetByID(ctx context.Context, id string) (*store.Expense, error)
	List(ctx context.Context, userID, visibility, from, to string) ([]store.Expense, error)
	Update(ctx context.Context, id, userID, description, category, visibility string, amount, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*store.Expense, error)
	Delete(ctx context.Context, id, userID string) error
	Settle(ctx context.Context, id, userID string) (*store.Expense, error)
}

// ─── ExpenseHandler ─────────────────────────────────────────────────

// ExpenseHandler handles expense endpoints.
type ExpenseHandler struct {
	expenses ExpenseRepository
}

// NewExpenseHandler creates a new ExpenseHandler.
func NewExpenseHandler(expenses ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{expenses: expenses}
}

// ─── Request/Response types ─────────────────────────────────────────

type createExpenseRequest struct {
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	Visibility      string  `json:"visibility"`
	SplitPercentage float64 `json:"split_percentage"`
	ExpenseDate     string  `json:"expense_date"`
	IsRecurring     bool    `json:"is_recurring"`
}

type updateExpenseRequest struct {
	Amount          float64 `json:"amount"`
	Description     string  `json:"description"`
	Category        string  `json:"category"`
	Visibility      string  `json:"visibility"`
	SplitPercentage float64 `json:"split_percentage"`
	ExpenseDate     string  `json:"expense_date"`
	IsRecurring     bool    `json:"is_recurring"`
}

// ─── Handlers ───────────────────────────────────────────────────────

// Create handles POST /api/v1/expenses.
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Description == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "description and positive amount are required")
		return
	}

	visibility := normalizeVisibility(req.Visibility, "private")
	if visibility != "private" && visibility != "shared" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "visibility must be private or shared")
		return
	}

	expenseDate, err := parseDate(req.ExpenseDate, time.Now())
	if err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid expense date")
		return
	}

	splitPercentage := req.SplitPercentage
	if splitPercentage <= 0 || splitPercentage > 100 {
		splitPercentage = 50
	}

	expense, err := h.expenses.Create(r.Context(), userID, req.Amount, req.Description, req.Category, visibility, splitPercentage, expenseDate, req.IsRecurring)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create expense")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": expense})
}

// List handles GET /api/v1/expenses.
func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	visibility := r.URL.Query().Get("visibility")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if visibility != "" && visibility != "private" && visibility != "shared" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid visibility filter")
		return
	}

	expenses, err := h.expenses.List(r.Context(), userID, visibility, from, to)
	if err != nil {
		log.Printf("failed to list expenses: %v", err)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list expenses")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expenses})
}

// Get handles GET /api/v1/expenses/{id}.
func (h *ExpenseHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "expense id is required")
		return
	}

	expense, err := h.expenses.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "expense not found")
		return
	}

	if expense.UserID != userID && expense.Visibility != "shared" {
		respondError(w, http.StatusForbidden, "FORBIDDEN", "not allowed to view this expense")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

// Update handles PUT /api/v1/expenses/{id}.
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "expense id is required")
		return
	}

	var req updateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Description == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "description and positive amount are required")
		return
	}

	visibility := normalizeVisibility(req.Visibility, "private")
	if visibility != "private" && visibility != "shared" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "visibility must be private or shared")
		return
	}

	expenseDate, err := parseDate(req.ExpenseDate, time.Now())
	if err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid expense date")
		return
	}

	splitPercentage := req.SplitPercentage
	if splitPercentage <= 0 || splitPercentage > 100 {
		splitPercentage = 50
	}

	expense, err := h.expenses.Update(r.Context(), id, userID, req.Description, req.Category, visibility, req.Amount, splitPercentage, expenseDate, req.IsRecurring)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "expense not found or not owned by user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

// Settle handles PATCH /api/v1/expenses/{id}/settle.
func (h *ExpenseHandler) Settle(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "expense id is required")
		return
	}

	expense, err := h.expenses.Settle(r.Context(), id, userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "expense not found or not shared")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": expense})
}

// Delete handles DELETE /api/v1/expenses/{id}.
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "expense id is required")
		return
	}

	if err := h.expenses.Delete(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "expense not found or not owned by user")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// normalizeVisibility returns a valid visibility value.
func normalizeVisibility(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// parseDate parses a date string or returns the default.
func parseDate(value string, defaultValue time.Time) (time.Time, error) {
	if value == "" {
		return defaultValue, nil
	}
	return time.Parse("2006-01-02", value)
}

// userIDFromContext extracts the user ID from request context.
func userIDFromContext(r *http.Request) (string, bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return "", false
	}
	return claims.UserID, true
}
