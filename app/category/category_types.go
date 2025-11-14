package category

import (
	"errors"

	"github.com/mytheresa/go-hiring-challenge/models"
)

var (
	ErrCategoryCodeRequired = errors.New("category code is required")
	ErrCategoryNameRequired = errors.New("category name is required")
	ErrCategoryCodeTooLong  = errors.New("category code must not exceed 32 characters")
	ErrCategoryNameTooLong  = errors.New("category name must not exceed 256 characters")
)

type CategoriesReader interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category *models.Category) error
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type CategoriesListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
