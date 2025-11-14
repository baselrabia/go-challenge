package category

import (
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoriesService struct {
	repo CategoriesReader
}

func NewCategoriesService(repo CategoriesReader) *CategoriesService {
	return &CategoriesService{
		repo: repo,
	}
}

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

func (s *CategoriesService) CreateCategory(req CreateCategoryRequest) (*CategoryResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

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

func (s *CategoriesService) validateCreateRequest(req CreateCategoryRequest) error {
	if strings.TrimSpace(req.Code) == "" {
		return ErrCategoryCodeRequired
	}

	if strings.TrimSpace(req.Name) == "" {
		return ErrCategoryNameRequired
	}

	if len(req.Code) > 32 {
		return ErrCategoryCodeTooLong
	}

	if len(req.Name) > 256 {
		return ErrCategoryNameTooLong
	}

	return nil
}
