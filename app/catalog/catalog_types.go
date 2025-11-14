package catalog

import "github.com/mytheresa/go-hiring-challenge/models"

// ProductsReader interface for fetching products
// This interface allows the handler to depend on behavior rather than concrete implementation
type ProductsReader interface {
	GetAllProducts() ([]models.Product, error)
	GetProductsWithPagination(offset, limit int, category string, priceLessThan *float64) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
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

type ProductDetail struct {
	Code     string          `json:"code"`
	Price    float64         `json:"price"`
	Category *Category       `json:"category,omitempty"`
	Variants []VariantDetail `json:"variants"`
}

type VariantDetail struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}
