// Package handlers provides HTTP handlers for products.
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/gentleman-programming/home-manager/internal/store"
)

// ─── ProductRepository ──────────────────────────────────────────────

// ProductRepository defines the interface for product data access.
type ProductRepository interface {
	Create(ctx context.Context, name, unit, category string) (*store.Product, error)
	GetByID(ctx context.Context, id string) (*store.Product, error)
	List(ctx context.Context, search string, limit, offset int) ([]store.Product, int, error)
	Update(ctx context.Context, id, name, unit, category string) (*store.Product, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, name string) (bool, error)
}

// PriceRepository defines the interface for price data access.
type PriceRepository interface {
	Create(ctx context.Context, productID, store string, amount float64) (*store.ProductPrice, error)
	ListByProduct(ctx context.Context, productID string, limit int) ([]store.ProductPrice, error)
}

// ─── ProductHandler ─────────────────────────────────────────────────

// ProductHandler handles product endpoints.
type ProductHandler struct {
	products ProductRepository
	prices   PriceRepository
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(products ProductRepository, prices PriceRepository) *ProductHandler {
	return &ProductHandler{
		products: products,
		prices:   prices,
	}
}

// ─── Request/Response types ─────────────────────────────────────────

type createProductRequest struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Category string `json:"category"`
}

type updateProductRequest struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Category string `json:"category"`
}

type createPriceRequest struct {
	Store  string  `json:"store"`
	Amount float64 `json:"amount"`
}

// ─── Handlers ───────────────────────────────────────────────────────

// Create handles POST /api/v1/products.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" || req.Unit == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "name and unit are required")
		return
	}

	exists, err := h.products.Exists(r.Context(), req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to check product")
		return
	}
	if exists {
		respondError(w, http.StatusConflict, "DUPLICATE_PRODUCT", "product already exists")
		return
	}

	product, err := h.products.Create(r.Context(), req.Name, req.Unit, req.Category)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create product")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": product})
}

// List handles GET /api/v1/products.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	limit := parseInt(r.URL.Query().Get("limit"), 20)
	offset := parseInt(r.URL.Query().Get("offset"), 0)

	products, total, err := h.products.List(r.Context(), search, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list products")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data": products,
		"meta": map[string]interface{}{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// Get handles GET /api/v1/products/{id}.
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product id is required")
		return
	}

	product, err := h.products.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": product})
}

// Update handles PUT /api/v1/products/{id}.
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product id is required")
		return
	}

	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Name == "" || req.Unit == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "name and unit are required")
		return
	}

	product, err := h.products.Update(r.Context(), id, req.Name, req.Unit, req.Category)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": product})
}

// Delete handles DELETE /api/v1/products/{id}.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product id is required")
		return
	}

	if err := h.products.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "product not found")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// CreatePrice handles POST /api/v1/products/{id}/prices.
func (h *ProductHandler) CreatePrice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product id is required")
		return
	}

	var req createPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	if req.Store == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "store and positive amount are required")
		return
	}

	price, err := h.prices.Create(r.Context(), id, req.Store, req.Amount)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create price")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{"data": price})
}

// ListPrices handles GET /api/v1/products/{id}/prices.
func (h *ProductHandler) ListPrices(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "BAD_REQUEST", "product id is required")
		return
	}

	limit := parseInt(r.URL.Query().Get("limit"), 10)

	prices, err := h.prices.ListByProduct(r.Context(), id, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to list prices")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": prices})
}

// parseInt parses an integer from a string with a default value.
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return defaultValue
	}
	return n
}
