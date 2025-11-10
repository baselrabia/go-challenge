package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// ProductsReader interface for fetching products
// This interface allows the handler to depend on behavior rather than concrete implementation
type ProductsReader interface {
	GetAllProducts() ([]models.Product, error)
	GetProductsWithPagination(offset, limit int, category string, priceLessThan *float64) ([]models.Product, int64, error)
}

type Response struct {
	Products []Product `json:"products"`
}

type PaginatedResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Offset   int       `json:"offset"`
	Limit    int       `json:"limit"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category *Category `json:"category,omitempty"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CatalogHandler struct {
	repo ProductsReader
}

func NewCatalogHandler(r ProductsReader) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	// Parse and validate pagination parameters
	offset := parseIntParam(r.URL.Query().Get("offset"), 0)
	if offset < 0 {
		offset = 0
	}

	limit := parseIntParam(r.URL.Query().Get("limit"), 10)
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

	// Get paginated products with filters
	res, total, err := h.repo.GetProductsWithPagination(offset, limit, category, priceLessThan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := make([]Product, len(res))
	for i, p := range res {
		product := Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
		}

		// Include category if available
		if p.Category != nil {
			product.Category = &Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			}
		}

		products[i] = product
	}

	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	response := PaginatedResponse{
		Products: products,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// parseIntParam parses a string parameter to int, returning defaultValue if parsing fails
func parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}
	var value int
	if _, err := fmt.Sscanf(param, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}
