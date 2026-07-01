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

func (m *mockDashboardRepository) PartnerBalance(ctx context.Context, userID string) (*store.PartnerBalance, error) {
	return m.balance, m.err
}

func TestDashboardHandler_Balance_Success(t *testing.T) {
	repo := &mockDashboardRepository{
		balance: &store.PartnerBalance{
			PartnerID:    "partner-1",
			PartnerEmail: "mariajose@home.local",
			Amount:       91,
			YouOwe:       false,
			YouAreOwed:   300,
			YouOweAmount: 209,
		},
	}
	handler := NewDashboardHandler(repo)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/balance", nil).WithContext(ctx)
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
	if data["amount"] != float64(91) {
		t.Errorf("amount = %v, want %v", data["amount"], 91)
	}
	if data["you_are_owed"] != float64(300) {
		t.Errorf("you_are_owed = %v, want %v", data["you_are_owed"], 300)
	}
	if data["you_owe_amount"] != float64(209) {
		t.Errorf("you_owe_amount = %v, want %v", data["you_owe_amount"], 209)
	}
	if data["you_owe"] != false {
		t.Errorf("you_owe = %v, want false", data["you_owe"])
	}
}

func TestDashboardHandler_Balance_Unauthorized(t *testing.T) {
	repo := &mockDashboardRepository{}
	handler := NewDashboardHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/balance", nil)
	rec := httptest.NewRecorder()

	handler.Balance(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
