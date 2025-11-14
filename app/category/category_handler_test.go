package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockCategoriesRepo struct {
	getAllFn func() ([]models.Category, error)
	createFn func(*models.Category) error
}

func (m *mockCategoriesRepo) GetAllCategories() ([]models.Category, error) {
	if m.getAllFn != nil {
		return m.getAllFn()
	}
	return nil, errors.New("not implemented")
}

func (m *mockCategoriesRepo) CreateCategory(category *models.Category) error {
	if m.createFn != nil {
		return m.createFn(category)
	}
	return errors.New("not implemented")
}

func TestHandleList_Success(t *testing.T) {
	categories := []models.Category{
		{ID: 1, Code: "CLOTHING", Name: "Clothing"},
		{ID: 2, Code: "SHOES", Name: "Shoes"},
		{ID: 3, Code: "ACCESSORIES", Name: "Accessories"},
	}

	repo := &mockCategoriesRepo{
		getAllFn: func() ([]models.Category, error) {
			return categories, nil
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleList(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp CategoriesListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	require.Len(t, resp.Categories, 3)
	assert.Equal(t, "CLOTHING", resp.Categories[0].Code)
	assert.Equal(t, "Clothing", resp.Categories[0].Name)
	assert.Equal(t, "SHOES", resp.Categories[1].Code)
	assert.Equal(t, "Shoes", resp.Categories[1].Name)
	assert.Equal(t, "ACCESSORIES", resp.Categories[2].Code)
	assert.Equal(t, "Accessories", resp.Categories[2].Name)
}

func TestHandleList_EmptyList(t *testing.T) {
	repo := &mockCategoriesRepo{
		getAllFn: func() ([]models.Category, error) {
			return []models.Category{}, nil
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	handler.HandleList(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp CategoriesListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Empty(t, resp.Categories)
}

func TestHandleList_DatabaseError(t *testing.T) {
	repo := &mockCategoriesRepo{
		getAllFn: func() ([]models.Category, error) {
			return nil, errors.New("database connection failed")
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))
	w := httptest.NewRecorder()

	handler.HandleList(w, httptest.NewRequest("GET", "/categories", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleCreate_Success(t *testing.T) {
	repo := &mockCategoriesRepo{
		createFn: func(category *models.Category) error {
			category.ID = 1
			return nil
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "ELECTRONICS", Name: "Electronics"})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var resp CategoryResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "ELECTRONICS", resp.Code)
	assert.Equal(t, "Electronics", resp.Name)
}

func TestHandleCreate_InvalidJSON(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	req := httptest.NewRequest("POST", "/categories", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_MissingCode(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "", Name: "Electronics"})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_MissingName(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "ELECTRONICS", Name: ""})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_CodeTooLong(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{
		Code: "THISISTOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOLONG",
		Name: "Test Category",
	})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_UppercaseConversion(t *testing.T) {
	repo := &mockCategoriesRepo{
		createFn: func(category *models.Category) error {
			// Business rule: code must be uppercased before storage
			assert.Equal(t, "ELECTRONICS", category.Code)
			category.ID = 1
			return nil
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "electronics", Name: "Electronics"})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp CategoryResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ELECTRONICS", resp.Code)
}

func TestHandleCreate_NameTooLong(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{
		Code: "TECH",
		Name: strings.Repeat("A", 257),
	})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_WhitespaceCode(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "   ", Name: "Electronics"})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_WhitespaceName(t *testing.T) {
	handler := NewCategoriesHandler(NewCategoriesService(&mockCategoriesRepo{}))

	body, _ := json.Marshal(CreateCategoryRequest{Code: "TECH", Name: "   "})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreate_DatabaseError(t *testing.T) {
	repo := &mockCategoriesRepo{
		createFn: func(category *models.Category) error {
			return errors.New("database constraint violation")
		},
	}

	handler := NewCategoriesHandler(NewCategoriesService(repo))
	body, _ := json.Marshal(CreateCategoryRequest{Code: "TECH", Name: "Technology"})
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreate(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
