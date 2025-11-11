package category

import (
	"errors"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoriesService contains business logic for category operations
type CategoriesService struct {
	repo CategoriesReader
}

// NewCategoriesService creates a new categories service
func NewCategoriesService(repo CategoriesReader) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}

// ListCategories retrieves all categories
func (s *CategoriesService) ListCategories() (*CategoriesListResponse, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		return nil, err
	}

	categoryDTOs := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		categoryDTOs[i] = CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	return &CategoriesListResponse{
		Categories: categoryDTOs,
	}, nil
}

// CreateCategory creates a new category with validation
func (s *CategoriesService) CreateCategory(req CreateCategoryRequest) (*CategoryResponse, error) {
	// Validate input
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Create category
	category := &models.Category{
		Code: strings.ToUpper(req.Code),
		Name: req.Name,
	}

	if err := s.repo.CreateCategory(category); err != nil {
		return nil, err
	}

	return &CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}, nil
}

// validateCreateRequest validates the create category request
func (s *CategoriesService) validateCreateRequest(req CreateCategoryRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("category code is required")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New("category name is required")
	}

	if len(req.Code) > 32 {
		return errors.New("category code must not exceed 32 characters")
	}

	if len(req.Name) > 256 {
		return errors.New("category name must not exceed 256 characters")
	}

	return nil
}
