package category

import (
	"encoding/json"
	"net/http"
)

// CategoriesHandler handles HTTP requests for the categories endpoints
type CategoriesHandler struct {
	service *CategoriesService
}

// NewCategoriesHandler creates a new categories handler
func NewCategoriesHandler(service *CategoriesService) *CategoriesHandler {
	return &CategoriesHandler{
		service: service,
	}
}

// HandleList handles GET /categories
func (h *CategoriesHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	// Call service to get categories
	response, err := h.service.ListCategories()
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

// HandleCreate handles POST /categories
func (h *CategoriesHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call service to create category
	response, err := h.service.CreateCategory(req)
	if err != nil {
		// Check for validation errors
		if err.Error() == "category code is required" ||
			err.Error() == "category name is required" ||
			err.Error() == "category code must not exceed 32 characters" ||
			err.Error() == "category name must not exceed 256 characters" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response with 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
