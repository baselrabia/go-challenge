package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProductsRepo struct {
	getByCodeFn func(string) (*models.Product, error)
}

func (m *mockProductsRepo) GetAllProducts() ([]models.Product, error) {
	return nil, nil
}

func (m *mockProductsRepo) GetProductsWithPagination(offset, limit int, category string, priceLessThan *float64) ([]models.Product, int64, error) {
	return nil, 0, nil
}

func (m *mockProductsRepo) GetProductByCode(code string) (*models.Product, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(code)
	}
	return nil, errors.New("not implemented")
}

func TestHandleGetByCode_Success(t *testing.T) {
	category := &models.Category{ID: 1, Code: "CLOTHING", Name: "Clothing"}
	product := &models.Product{
		ID:         1,
		Code:       "PROD001",
		Price:      decimal.NewFromFloat(10.99),
		CategoryID: &category.ID,
		Category:   category,
		Variants: []models.Variant{
			{ID: 1, ProductID: 1, Name: "Red", SKU: "SKU001-R", Price: decimal.NewFromFloat(11.99)},
			{ID: 2, ProductID: 1, Name: "Blue", SKU: "SKU001-B", Price: decimal.Zero},
			{ID: 3, ProductID: 1, Name: "Green", SKU: "SKU001-G", Price: decimal.Zero},
		},
	}

	repo := &mockProductsRepo{
		getByCodeFn: func(code string) (*models.Product, error) {
			if code == "PROD001" {
				return product, nil
			}
			return nil, errors.New("not found")
		},
	}

	handler := NewCatalogHandler(NewCatalogService(repo))
	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp ProductDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.Equal(t, "PROD001", resp.Code)
	assert.Equal(t, 10.99, resp.Price)
	assert.Equal(t, "CLOTHING", resp.Category.Code)
	assert.Equal(t, "Clothing", resp.Category.Name)

	require.Len(t, resp.Variants, 3)

	// Variant with explicit price keeps it
	assert.Equal(t, 11.99, resp.Variants[0].Price)

	// Variants with zero price inherit from product
	assert.Equal(t, 10.99, resp.Variants[1].Price)
	assert.Equal(t, 10.99, resp.Variants[2].Price)
}

func TestHandleGetByCode_NotFound(t *testing.T) {
	repo := &mockProductsRepo{
		getByCodeFn: func(code string) (*models.Product, error) {
			return nil, errors.New("not found")
		},
	}

	handler := NewCatalogHandler(NewCatalogService(repo))
	req := httptest.NewRequest("GET", "/catalog/INVALID", nil)
	req.SetPathValue("code", "INVALID")
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleGetByCode_EmptyCode(t *testing.T) {
	handler := NewCatalogHandler(NewCatalogService(&mockProductsRepo{}))
	req := httptest.NewRequest("GET", "/catalog/", nil)
	req.SetPathValue("code", "")
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGetByCode_NoCategory(t *testing.T) {
	product := &models.Product{
		ID:    1,
		Code:  "NOCATEGORY",
		Price: decimal.NewFromFloat(25.00),
		Variants: []models.Variant{
			{ID: 1, ProductID: 1, Name: "Default", SKU: "NOCAT-1", Price: decimal.NewFromFloat(26.00)},
		},
	}

	repo := &mockProductsRepo{
		getByCodeFn: func(code string) (*models.Product, error) {
			return product, nil
		},
	}

	handler := NewCatalogHandler(NewCatalogService(repo))
	req := httptest.NewRequest("GET", "/catalog/NOCATEGORY", nil)
	req.SetPathValue("code", "NOCATEGORY")
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp ProductDetail
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	assert.Equal(t, "NOCATEGORY", resp.Code)
	assert.Nil(t, resp.Category)
	assert.Len(t, resp.Variants, 1)
}

func TestHandleGetByCode_NoVariants(t *testing.T) {
	category := &models.Category{ID: 1, Code: "SHOES", Name: "Shoes"}
	product := &models.Product{
		ID:         1,
		Code:       "SIMPLE",
		Price:      decimal.NewFromFloat(50.00),
		CategoryID: &category.ID,
		Category:   category,
		Variants:   []models.Variant{},
	}

	repo := &mockProductsRepo{
		getByCodeFn: func(code string) (*models.Product, error) {
			return product, nil
		},
	}

	handler := NewCatalogHandler(NewCatalogService(repo))
	req := httptest.NewRequest("GET", "/catalog/SIMPLE", nil)
	req.SetPathValue("code", "SIMPLE")
	w := httptest.NewRecorder()

	handler.HandleGetByCode(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp ProductDetail
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "SIMPLE", resp.Code)
	assert.Equal(t, 50.00, resp.Price)
	assert.Empty(t, resp.Variants)
}
