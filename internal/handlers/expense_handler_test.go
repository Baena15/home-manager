package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"
	"github.com/gentleman-programming/home-manager/pkg/auth"
)

type mockExpenseRepository struct {
	expense  *store.Expense
	expenses []store.Expense
	err      error
}

func (m *mockExpenseRepository) Create(ctx context.Context, userID string, amount float64, description, category, visibility string, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*store.Expense, error) {
	return m.expense, m.err
}

func (m *mockExpenseRepository) GetByID(ctx context.Context, id string) (*store.Expense, error) {
	return m.expense, m.err
}

func (m *mockExpenseRepository) List(ctx context.Context, userID, visibility, from, to string) ([]store.Expense, error) {
	return m.expenses, m.err
}

func (m *mockExpenseRepository) Update(ctx context.Context, id, userID, description, category, visibility string, amount, splitPercentage float64, expenseDate time.Time, isRecurring bool) (*store.Expense, error) {
	return m.expense, m.err
}

func (m *mockExpenseRepository) Delete(ctx context.Context, id, userID string) error {
	return m.err
}

func (m *mockExpenseRepository) Settle(ctx context.Context, id, userID string) (*store.Expense, error) {
	return m.expense, m.err
}

func TestExpenseHandler_Settle_Success(t *testing.T) {
	settledAt := time.Now()
	repo := &mockExpenseRepository{
		expense: &store.Expense{
			ID:              "exp-1",
			UserID:          "user-2",
			Amount:          600,
			Description:     "crédito",
			Visibility:      "shared",
			SplitPercentage: 50,
			SettledAt:       &settledAt,
		},
	}
	handler := NewExpenseHandler(repo)

	r := chi.NewRouter()
	r.Patch("/api/v1/expenses/{id}/settle", handler.Settle)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/expenses/exp-1/settle", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestExpenseHandler_Settle_NotFound(t *testing.T) {
	repo := &mockExpenseRepository{
		err: errors.New("not found"),
	}
	handler := NewExpenseHandler(repo)

	r := chi.NewRouter()
	r.Patch("/api/v1/expenses/{id}/settle", handler.Settle)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/expenses/exp-1/settle", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
