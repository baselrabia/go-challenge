package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// ProductsReader interface for fetching products
// This interface allows the handler to depend on behavior rather than concrete implementation
type ProductsReader interface {
	GetAllProducts() ([]models.Product, error)
}

type Response struct {
	Products []Product `json:"products"`
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
	res, err := h.repo.GetAllProducts()
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

	response := Response{
		Products: products,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
