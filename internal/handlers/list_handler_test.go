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

type mockListRepository struct {
	list  *store.ShoppingList
	lists []store.ShoppingList
	total int
	err   error
}

func (m *mockListRepository) Create(ctx context.Context, name, createdBy string) (*store.ShoppingList, error) {
	return m.list, m.err
}

func (m *mockListRepository) GetByID(ctx context.Context, id string) (*store.ShoppingList, error) {
	return m.list, m.err
}

func (m *mockListRepository) List(ctx context.Context, limit, offset int) ([]store.ShoppingList, int, error) {
	return m.lists, m.total, m.err
}

func (m *mockListRepository) Update(ctx context.Context, id, name string) (*store.ShoppingList, error) {
	return m.list, m.err
}

func (m *mockListRepository) UpdateStatus(ctx context.Context, id, status string) (*store.ShoppingList, error) {
	if m.list != nil {
		m.list.Status = status
	}
	return m.list, m.err
}

func (m *mockListRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

type mockListItemRepository struct {
	item  *store.ShoppingListItem
	items []store.ShoppingListItem
	err   error
}

func (m *mockListItemRepository) Create(ctx context.Context, listID, productID string, quantity, unitPrice float64) (*store.ShoppingListItem, error) {
	return m.item, m.err
}

func (m *mockListItemRepository) ListByList(ctx context.Context, listID string) ([]store.ShoppingListItem, error) {
	return m.items, m.err
}

func (m *mockListItemRepository) UpdatePurchased(ctx context.Context, id string, purchased bool) (*store.ShoppingListItem, error) {
	if m.item != nil {
		m.item.Purchased = purchased
	}
	return m.item, m.err
}

func (m *mockListItemRepository) Delete(ctx context.Context, id string) error {
	return m.err
}

func setupListRouter(h *ListHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/v1/lists", h.Create)
	r.Get("/api/v1/lists", h.List)
	r.Get("/api/v1/lists/{id}", h.Get)
	r.Put("/api/v1/lists/{id}", h.Update)
	r.Delete("/api/v1/lists/{id}", h.Delete)
	r.Patch("/api/v1/lists/{id}/status", h.UpdateStatus)
	r.Post("/api/v1/lists/{id}/items", h.AddItem)
	r.Patch("/api/v1/lists/{id}/items/{item_id}", h.UpdateItem)
	r.Delete("/api/v1/lists/{id}/items/{item_id}", h.RemoveItem)
	return r
}

func TestListHandler_Create_Success(t *testing.T) {
	lists := &mockListRepository{
		list: &store.ShoppingList{
			ID:        "list-1",
			Name:      "Compra semanal",
			Status:    "active",
			CreatedBy: "user-1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	items := &mockListItemRepository{}
	products := &mockProductRepository{}
	handler := NewListHandler(lists, items, products)
	router := setupListRouter(handler)

	ctx := context.WithValue(context.Background(), middleware.ClaimsKey, &auth.Claims{UserID: "user-1"})
	body, _ := json.Marshal(map[string]string{"name": "Compra semanal"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lists", bytes.NewReader(body)).WithContext(ctx)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestListHandler_List(t *testing.T) {
	lists := &mockListRepository{
		lists: []store.ShoppingList{
			{ID: "list-1", Name: "Compra semanal", Status: "active"},
		},
		total: 1,
	}
	items := &mockListItemRepository{}
	products := &mockProductRepository{}
	handler := NewListHandler(lists, items, products)
	router := setupListRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lists", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestListHandler_AddItem_WithLatestPrice(t *testing.T) {
	latest := 4.50
	lists := &mockListRepository{
		list: &store.ShoppingList{ID: "list-1"},
	}
	items := &mockListItemRepository{
		item: &store.ShoppingListItem{
			ID:        "item-1",
			ListID:    "list-1",
			ProductID: "prod-1",
			Quantity:  2,
			UnitPrice: 4.50,
			Total:     9.00,
		},
	}
	products := &mockProductRepository{
		product: &store.Product{
			ID:          "prod-1",
			Name:        "Carne picada",
			LatestPrice: &latest,
		},
	}
	handler := NewListHandler(lists, items, products)
	router := setupListRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{"product_id": "prod-1", "quantity": 2})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lists/list-1/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
}

func TestListHandler_AddItem_WithoutPrice(t *testing.T) {
	lists := &mockListRepository{
		list: &store.ShoppingList{ID: "list-1"},
	}
	items := &mockListItemRepository{}
	products := &mockProductRepository{
		product: &store.Product{
			ID:   "prod-1",
			Name: "Carne picada",
		},
	}
	handler := NewListHandler(lists, items, products)
	router := setupListRouter(handler)

	body, _ := json.Marshal(map[string]interface{}{"product_id": "prod-1", "quantity": 2})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/lists/list-1/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestListHandler_UpdateItem(t *testing.T) {
	lists := &mockListRepository{}
	items := &mockListItemRepository{
		item: &store.ShoppingListItem{
			ID:        "item-1",
			Purchased: false,
		},
	}
	products := &mockProductRepository{}
	handler := NewListHandler(lists, items, products)
	router := setupListRouter(handler)

	body, _ := json.Marshal(map[string]bool{"purchased": true})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/lists/list-1/items/item-1", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
