// Package handlers provides HTTP handlers for settlements.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── SettlementRepository ───────────────────────────────────────────

// SettlementRepository defines the interface for settlement data access.
type SettlementRepository interface {
	Create(ctx context.Context, fromUserID, toUserID string, amount float64, description string, settlementDate time.Time) (*store.Settlement, error)
	List(ctx context.Context, userID, from, to string) ([]store.Settlement, error)
	GetByID(ctx context.Context, id string) (*store.Settlement, error)
	Delete(ctx context.Context, id, userID string) error
}

// ─── SettlementHandler ──────────────────────────────────────────────

// SettlementHandler handles settlement endpoints.
type SettlementHandler struct {
	settlements SettlementRepository
	users       UserRepository
}

// NewSettlementHandler creates a new SettlementHandler.
func NewSettlementHandler(settlements SettlementRepository, users UserRepository) *SettlementHandler {
	return &SettlementHandler{
		settlements: settlements,
		users:       users,
	}
}

// ─── Request/Response types ─────────────────────────────────────────

type createSettlementRequest struct {
	Amount         float64 `json:"amount"`
	Description    string  `json:"description"`
	ToUserID       string  `json:"to_user_id"`
	SettlementDate string  `json:"settlement_date"`
}

// ─── Handlers ───────────────────────────────────────────────────────

// Create handles POST /api/v1/settlements.
func (h *SettlementHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	var req createSettlementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Amount <= 0 || req.ToUserID == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "amount and receiver are required")
		return
	}

	if req.ToUserID == userID {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "cannot send payment to yourself")
		return
	}

	if _, err := h.users.GetByID(r.Context(), req.ToUserID); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "receiver not found")
		return
	}

	settlementDate, err := parseDate(req.SettlementDate, time.Now())
	if err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid settlement date")
		return
	}

	settlement, err := h.settlements.Create(r.Context(), userID, req.ToUserID, req.Amount, req.Description, settlementDate)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create settlement")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": settlement})
}

// List handles GET /api/v1/settlements.
func (h *SettlementHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	settlements, err := h.settlements.List(r.Context(), userID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list settlements")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": settlements})
}

// Get handles GET /api/v1/settlements/{id}.
func (h *SettlementHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "settlement id is required")
		return
	}

	settlement, err := h.settlements.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "settlement not found")
		return
	}

	if settlement.FromUserID != userID && settlement.ToUserID != userID {
		respondError(w, http.StatusForbidden, "FORBIDDEN", "not allowed to view this settlement")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": settlement})
}

// Delete handles DELETE /api/v1/settlements/{id}.
func (h *SettlementHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "settlement id is required")
		return
	}

	if err := h.settlements.Delete(r.Context(), id, userID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "settlement not found or not involved user")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
