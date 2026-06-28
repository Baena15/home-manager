// Package handlers provides HTTP handlers for incomes.
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── IncomeRepository ───────────────────────────────────────────────

// IncomeRepository defines the interface for income data access.
type IncomeRepository interface {
	Create(ctx context.Context, userID string, amount float64, description, category, visibility string, incomeDate time.Time, isRecurring bool) (*store.Income, error)
	GetByID(ctx context.Context, id string) (*store.Income, error)
	List(ctx context.Context, userID, visibility, from, to string) ([]store.Income, error)
	Update(ctx context.Context, id, userID, description, category, visibility string, amount float64, incomeDate time.Time, isRecurring bool) (*store.Income, error)
	Delete(ctx context.Context, id, userID string) error
}

// ─── IncomeHandler ──────────────────────────────────────────────────

// IncomeHandler handles income endpoints.
type IncomeHandler struct {
	incomes IncomeRepository
}

// NewIncomeHandler creates a new IncomeHandler.
func NewIncomeHandler(incomes IncomeRepository) *IncomeHandler {
	return &IncomeHandler{incomes: incomes}
}

// ─── Request/Response types ─────────────────────────────────────────

type createIncomeRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Visibility  string  `json:"visibility"`
	IncomeDate  string  `json:"income_date"`
	IsRecurring bool    `json:"is_recurring"`
}

type updateIncomeRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Visibility  string  `json:"visibility"`
	IncomeDate  string  `json:"income_date"`
	IsRecurring bool    `json:"is_recurring"`
}

// ─── Handlers ───────────────────────────────────────────────────────

// Create handles POST /api/v1/incomes.
func (h *IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req createIncomeRequest
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

	incomeDate, err := parseDate(req.IncomeDate, time.Now())
	if err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid income date")
		return
	}

	income, err := h.incomes.Create(r.Context(), userID, req.Amount, req.Description, req.Category, visibility, incomeDate, req.IsRecurring)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create income")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": income})
}

// List handles GET /api/v1/incomes.
func (h *IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
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

	incomes, err := h.incomes.List(r.Context(), userID, visibility, from, to)
	if err != nil {
		log.Printf("failed to list incomes: %v", err)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list incomes")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": incomes})
}

// Get handles GET /api/v1/incomes/{id}.
func (h *IncomeHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "income id is required")
		return
	}

	income, err := h.incomes.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "income not found")
		return
	}

	if income.UserID != userID && income.Visibility != "shared" {
		respondError(w, http.StatusForbidden, "FORBIDDEN", "not allowed to view this income")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": income})
}

// Update handles PUT /api/v1/incomes/{id}.
func (h *IncomeHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "income id is required")
		return
	}

	var req updateIncomeRequest
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

	incomeDate, err := parseDate(req.IncomeDate, time.Now())
	if err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid income date")
		return
	}

	income, err := h.incomes.Update(r.Context(), id, userID, req.Description, req.Category, visibility, req.Amount, incomeDate, req.IsRecurring)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "income not found or not owned by user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": income})
}

// Delete handles DELETE /api/v1/incomes/{id}.
func (h *IncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "income id is required")
		return
	}

	if err := h.incomes.Delete(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "income not found or not owned by user")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
