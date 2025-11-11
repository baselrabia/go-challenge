package category

import "github.com/mytheresa/go-hiring-challenge/models"

// CategoriesReader interface for fetching and creating categories
type CategoriesReader interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}

// CategoryResponse represents a category in the API response
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CategoriesListResponse represents the list of categories response
type CategoriesListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
