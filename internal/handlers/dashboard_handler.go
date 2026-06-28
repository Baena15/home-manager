// Package handlers provides HTTP handlers for the dashboard.
package handlers

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── DashboardRepository ────────────────────────────────────────────

// DashboardRepository defines the interface for dashboard aggregations.
type DashboardRepository interface {
	MonthlySummary(ctx context.Context, userID, month string) (*store.MonthlySummary, error)
	MonthlyTotals(ctx context.Context, userID, year string) ([]store.MonthData, error)
}

// ─── DashboardHandler ───────────────────────────────────────────────

// DashboardHandler handles dashboard endpoints.
type DashboardHandler struct {
	dashboard DashboardRepository
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(dashboard DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard}
}

// ─── Handlers ───────────────────────────────────────────────────────

// Summary handles GET /api/v1/dashboard/summary.
func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	month := r.URL.Query().Get("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	if !isValidMonth(month) {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid month format, expected YYYY-MM")
		return
	}

	summary, err := h.dashboard.MonthlySummary(r.Context(), userID, month)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to calculate summary")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": summary})
}

// Monthly handles GET /api/v1/dashboard/monthly.
func (h *DashboardHandler) Monthly(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
		return
	}

	year := r.URL.Query().Get("year")
	if year == "" {
		year = time.Now().Format("2006")
	}
	if !isValidYear(year) {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid year format, expected YYYY")
		return
	}

	data, err := h.dashboard.MonthlyTotals(r.Context(), userID, year)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to calculate monthly totals")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

// isValidMonth checks if a string matches YYYY-MM format.
func isValidMonth(value string) bool {
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, value)
	return matched
}

// isValidYear checks if a string matches YYYY format.
func isValidYear(value string) bool {
	matched, _ := regexp.MatchString(`^\d{4}$`, value)
	return matched
}
