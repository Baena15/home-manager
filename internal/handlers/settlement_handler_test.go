package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"
	"github.com/gentleman-programming/home-manager/pkg/auth"
)

type mockSettlementRepository struct {
	settlement  *store.Settlement
	settlements []store.Settlement
	err         error
}

func (m *mockSettlementRepository) Create(ctx context.Context, fromUserID, toUserID string, amount float64, description string, settlementDate time.Time) (*store.Settlement, error) {
	return m.settlement, m.err
}

func (m *mockSettlementRepository) List(ctx context.Context, userID, from, to string) ([]store.Settlement, error) {
	return m.settlements, m.err
}

func (m *mockSettlementRepository) GetByID(ctx context.Context, id string) (*store.Settlement, error) {
	return m.settlement, m.err
}

func (m *mockSettlementRepository) Delete(ctx context.Context, id, userID string) error {
	return m.err
}

func TestSettlementHandler_Create_Success(t *testing.T) {
	settlementRepo := &mockSettlementRepository{
		settlement: &store.Settlement{
			ID:             "set-1",
			FromUserID:     "user-1",
			ToUserID:       "user-2",
			Amount:         300,
			Description:    "Pago del crédito",
			SettlementDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	userRepo := &mockUserRepository{
		user: &store.User{
			ID:    "user-2",
			Email: "partner@home.local",
			Role:  "partner",
		},
	}
	handler := NewSettlementHandler(settlementRepo, userRepo)

	body, _ := json.Marshal(createSettlementRequest{
		Amount:         300,
		Description:    "Pago del crédito",
		ToUserID:       "user-2",
		SettlementDate: "2026-07-01",
	})

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/settlements", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestSettlementHandler_Create_SelfPayment(t *testing.T) {
	settlementRepo := &mockSettlementRepository{}
	userRepo := &mockUserRepository{}
	handler := NewSettlementHandler(settlementRepo, userRepo)

	body, _ := json.Marshal(createSettlementRequest{
		Amount:   100,
		ToUserID: "user-1",
	})

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/settlements", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestSettlementHandler_List(t *testing.T) {
	settlementRepo := &mockSettlementRepository{
		settlements: []store.Settlement{
			{
				ID:             "set-1",
				FromUserID:     "user-2",
				FromUserEmail:  "partner@home.local",
				ToUserID:       "user-1",
				ToUserEmail:    "owner@home.local",
				Amount:         300,
				SettlementDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	userRepo := &mockUserRepository{}
	handler := NewSettlementHandler(settlementRepo, userRepo)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/settlements?from=2026-07-01&to=2026-07-31", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	handler.List(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestSettlementHandler_Delete_Success(t *testing.T) {
	settlementRepo := &mockSettlementRepository{}
	userRepo := &mockUserRepository{}
	handler := NewSettlementHandler(settlementRepo, userRepo)

	r := chi.NewRouter()
	r.Delete("/api/v1/settlements/{id}", handler.Delete)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/settlements/set-1", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}
