package catalog

// CatalogService contains business logic for catalog operations
type CatalogService struct {
	repo ProductsReader
}

// NewCatalogService creates a new catalog service
func NewCatalogService(repo ProductsReader) *CatalogService {
	return &CatalogService{
		repo: repo,
	}
}

// ListProducts retrieves a paginated list of products with optional filters
func (s *CatalogService) ListProducts(offset, limit int, category string, priceLessThan *float64) (*PaginatedResponse, error) {
	// Get products from repository
	products, total, err := s.repo.GetProductsWithPagination(offset, limit, category, priceLessThan)
	if err != nil {
		return nil, err
	}

	// Map domain models to response DTOs
	productDTOs := make([]Product, len(products))
	for i, p := range products {
		productDTOs[i] = mapProductToDTO(p)
	}

	return &PaginatedResponse{
		Products: productDTOs,
		Total:    total,
		Offset:   offset,
		Limit:    limit,
	}, nil
}

// GetProductDetails retrieves detailed product information with variants
// Business Rule: Variants without specific prices inherit the product's price
func (s *CatalogService) GetProductDetails(code string) (*ProductDetail, error) {
	// Get product from repository
	product, err := s.repo.GetProductByCode(code)
	if err != nil {
		return nil, err
	}

	// Apply business logic: price inheritance
	return mapProductToDetailDTO(product), nil
}
