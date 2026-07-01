package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"
	"github.com/gentleman-programming/home-manager/pkg/auth"
)

type mockDashboardRepository struct {
	summary *store.MonthlySummary
	totals  []store.MonthData
	balance *store.PartnerBalance
	err     error
}

func (m *mockDashboardRepository) MonthlySummary(ctx context.Context, userID, month string) (*store.MonthlySummary, error) {
	return m.summary, m.err
}

func (m *mockDashboardRepository) MonthlyTotals(ctx context.Context, userID, year string) ([]store.MonthData, error) {
	return m.totals, m.err
}

func (m *mockDashboardRepository) PartnerBalance(ctx context.Context, userID, month string) (*store.PartnerBalance, error) {
	return m.balance, m.err
}

func TestDashboardHandler_Balance_Success(t *testing.T) {
	repo := &mockDashboardRepository{
		balance: &store.PartnerBalance{
			PartnerID:    "partner-1",
			PartnerEmail: "mariajose@home.local",
			Amount:       300,
			YouOwe:       false,
		},
	}
	handler := NewDashboardHandler(repo)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/balance?month=2026-07", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Balance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("response data is not an object")
	}

	if data["partner_email"] != "mariajose@home.local" {
		t.Errorf("partner_email = %q, want %q", data["partner_email"], "mariajose@home.local")
	}
	if data["amount"] != float64(300) {
		t.Errorf("amount = %v, want %v", data["amount"], 300)
	}
	if data["you_owe"] != false {
		t.Errorf("you_owe = %v, want false", data["you_owe"])
	}
}

func TestDashboardHandler_Balance_InvalidMonth(t *testing.T) {
	repo := &mockDashboardRepository{}
	handler := NewDashboardHandler(repo)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/balance?month=invalid", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Balance(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestDashboardHandler_Balance_Unauthorized(t *testing.T) {
	repo := &mockDashboardRepository{}
	handler := NewDashboardHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/balance?month=2026-07", nil)
	rec := httptest.NewRecorder()

	handler.Balance(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
