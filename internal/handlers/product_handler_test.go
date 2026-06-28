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

	"github.com/gentleman-programming/home-manager/internal/store"
)

type mockProductRepository struct {
	product  *store.Product
	products []store.Product
	total    int
	exists   bool
	err      error
}

func (m *mockProductRepository) Create(ctx context.Context, name, unit, category string) (*store.Product, error) {
	return m.product, m.err
}

func (m *mockProductRepository) GetByID(ctx context.Context, id string) (*store.Product, error) {
	return m.product, m.err
}

func (m *mockProductRepository) List(ctx context.Context, search string, limit, offset int) ([]store.Product, int, error) {
	return m.products, m.total, m.err
}

func (m *mockProductRepository) Update(ctx context.Context, id, name, unit, category string) (*store.Product, error) {
	return m.product, m.err
}

func (m *mockProductRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

func (m *mockProductRepository) Exists(ctx context.Context, name string) (bool, error) {
	return m.exists, m.err
}

type mockPriceRepository struct {
	price  *store.ProductPrice
	prices []store.ProductPrice
	err    error
}

func (m *mockPriceRepository) Create(ctx context.Context, productID, store string, amount float64) (*store.ProductPrice, error) {
	return m.price, m.err
}

func (m *mockPriceRepository) ListByProduct(ctx context.Context, productID string, limit int) ([]store.ProductPrice, error) {
	return m.prices, m.err
}

func setupProductRouter(h *ProductHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/v1/products", h.Create)
	r.Get("/api/v1/products", h.List)
	r.Get("/api/v1/products/{id}", h.Get)
	r.Put("/api/v1/products/{id}", h.Update)
	r.Delete("/api/v1/products/{id}", h.Delete)
	r.Post("/api/v1/products/{id}/prices", h.CreatePrice)
	r.Get("/api/v1/products/{id}/prices", h.ListPrices)
	return r
}

func TestProductHandler_Create_Success(t *testing.T) {
	repo := &mockProductRepository{
		product: &store.Product{
			ID:        "prod-1",
			Name:      "Carne picada",
			Unit:      "g",
			Category:  "carnicería",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	prices := &mockPriceRepository{}
	handler := NewProductHandler(repo, prices)
	router := setupProductRouter(handler)

	body, _ := json.Marshal(map[string]string{"name": "Carne picada", "unit": "g", "category": "carnicería"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestProductHandler_Create_Duplicate(t *testing.T) {
	repo := &mockProductRepository{
		exists: true,
	}
	prices := &mockPriceRepository{}
	handler := NewProductHandler(repo, prices)
	router := setupProductRouter(handler)

	body, _ := json.Marshal(map[string]string{"name": "Carne picada", "unit": "g"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestProductHandler_List(t *testing.T) {
	repo := &mockProductRepository{
		products: []store.Product{
			{ID: "prod-1", Name: "Carne picada", Unit: "g"},
		},
		total: 1,
	}
	prices := &mockPriceRepository{}
	handler := NewProductHandler(repo, prices)
	router := setupProductRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestProductHandler_CreatePrice(t *testing.T) {
	repo := &mockProductRepository{}
	prices := &mockPriceRepository{
		price: &store.ProductPrice{
			ID:         "price-1",
			ProductID:  "prod-1",
			Store:      "Mercadona",
			Amount:     4.50,
			RecordedAt: time.Now(),
		},
	}
	handler := NewProductHandler(repo, prices)
	router := setupProductRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{"store": "Mercadona", "amount": 4.50})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/prod-1/prices", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
}
