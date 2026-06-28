// Package handlers provides HTTP handlers for shopping lists.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/middleware"
	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── ListRepository ─────────────────────────────────────────────────

// ListRepository defines the interface for shopping list data access.
type ListRepository interface {
	Create(ctx context.Context, name, createdBy string) (*store.ShoppingList, error)
	GetByID(ctx context.Context, id string) (*store.ShoppingList, error)
	List(ctx context.Context, limit, offset int) ([]store.ShoppingList, int, error)
	Update(ctx context.Context, id, name string) (*store.ShoppingList, error)
	UpdateStatus(ctx context.Context, id, status string) (*store.ShoppingList, error)
	Delete(ctx context.Context, id string) error
}

// ListItemRepository defines the interface for shopping list item data access.
type ListItemRepository interface {
	Create(ctx context.Context, listID, productID string, quantity, unitPrice float64) (*store.ShoppingListItem, error)
	ListByList(ctx context.Context, listID string) ([]store.ShoppingListItem, error)
	UpdatePurchased(ctx context.Context, id string, purchased bool) (*store.ShoppingListItem, error)
	Delete(ctx context.Context, id string) error
}

// ─── ListHandler ────────────────────────────────────────────────────

// ListHandler handles shopping list endpoints.
type ListHandler struct {
	lists    ListRepository
	items    ListItemRepository
	products ProductRepository
}

// NewListHandler creates a new ListHandler.
func NewListHandler(lists ListRepository, items ListItemRepository, products ProductRepository) *ListHandler {
	return &ListHandler{
		lists:    lists,
		items:    items,
		products: products,
	}
}

// ─── Request/Response types ─────────────────────────────────────────

type createListRequest struct {
	Name string `json:"name"`
}

type updateListRequest struct {
	Name string `json:"name"`
}

type updateListStatusRequest struct {
	Status string `json:"status"`
}

type addItemRequest struct {
	ProductID   string   `json:"product_id"`
	Quantity    float64  `json:"quantity"`
	CustomPrice *float64 `json:"custom_price,omitempty"`
}

type updateItemRequest struct {
	Purchased bool `json:"purchased"`
}

// ─── Handlers ───────────────────────────────────────────────────────

// Create handles POST /api/v1/lists.
func (h *ListHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "name is required")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "missing claims")
		return
	}

	list, err := h.lists.Create(r.Context(), req.Name, claims.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create list")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": list})
}

// List handles GET /api/v1/lists.
func (h *ListHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := parseInt(r.URL.Query().Get("limit"), 20)
	offset := parseInt(r.URL.Query().Get("offset"), 0)

	lists, total, err := h.lists.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list lists")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": lists,
		"meta": map[string]interface{}{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// Get handles GET /api/v1/lists/{id}.
func (h *ListHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "list id is required")
		return
	}

	list, err := h.lists.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "list not found")
		return
	}

	items, err := h.items.ListByList(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list items")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"list":  list,
			"items": items,
		},
	})
}

// Update handles PUT /api/v1/lists/{id}.
func (h *ListHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "list id is required")
		return
	}

	var req updateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "name is required")
		return
	}

	list, err := h.lists.Update(r.Context(), id, req.Name)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "list not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// UpdateStatus handles PATCH /api/v1/lists/{id}/status.
func (h *ListHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "list id is required")
		return
	}

	var req updateListStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Status != "active" && req.Status != "completed" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "status must be active or completed")
		return
	}

	list, err := h.lists.UpdateStatus(r.Context(), id, req.Status)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "list not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

// Delete handles DELETE /api/v1/lists/{id}.
func (h *ListHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "list id is required")
		return
	}

	if err := h.lists.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "list not found")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// AddItem handles POST /api/v1/lists/{id}/items.
func (h *ListHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	listID := chi.URLParam(r, "id")
	if listID == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "list id is required")
		return
	}

	var req addItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.ProductID == "" || req.Quantity <= 0 {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product_id and positive quantity are required")
		return
	}

	var unitPrice float64
	if req.CustomPrice != nil && *req.CustomPrice > 0 {
		unitPrice = *req.CustomPrice
	} else {
		product, err := h.products.GetByID(r.Context(), req.ProductID)
		if err != nil {
			respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product not found")
			return
		}
		if product.LatestPrice == nil {
			respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product has no price, provide custom_price")
			return
		}
		unitPrice = *product.LatestPrice
	}

	item, err := h.items.Create(r.Context(), listID, req.ProductID, req.Quantity, unitPrice)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to add item")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": item})
}

// UpdateItem handles PATCH /api/v1/lists/{id}/items/{item_id}.
func (h *ListHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "item id is required")
		return
	}

	var req updateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	item, err := h.items.UpdatePurchased(r.Context(), itemID, req.Purchased)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "item not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": item})
}

// RemoveItem handles DELETE /api/v1/lists/{id}/items/{item_id}.
func (h *ListHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item_id")
	if itemID == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "item id is required")
		return
	}

	if err := h.items.Delete(r.Context(), itemID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "item not found")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}
