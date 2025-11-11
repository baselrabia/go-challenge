package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/utils"
)

// CatalogHandler handles HTTP requests for the catalog endpoints
type CatalogHandler struct {
	service *CatalogService
}

// NewCatalogHandler creates a new catalog handler
func NewCatalogHandler(service *CatalogService) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse and validate pagination parameters
	offset := utils.ParseIntParam(r.URL.Query().Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}

	limit := utils.ParseIntParam(r.URL.Query().Get("limit"), 10)
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	// Parse filter parameters
	category := r.URL.Query().Get("category")

	var priceLessThan *float64
	if priceLessThanStr := r.URL.Query().Get("priceLessThan"); priceLessThanStr != "" {
		var price float64
		if _, err := fmt.Sscanf(priceLessThanStr, "%f", &price); err == nil && price > 0 {
			priceLessThan = &price
		}
	}

	// Call service to get products
	response, err := h.service.ListProducts(offset, limit, category, priceLessThan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	// Extract product code from URL path
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Product code is required", http.StatusBadRequest)
		return
	}

	// Call service to get product details
	response, err := h.service.GetProductDetails(code)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
