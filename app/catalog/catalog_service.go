package catalog

type CatalogService struct {
	repo ProductsReader
}

func NewCatalogService(repo ProductsReader) *CatalogService {
	return &CatalogService{
		repo: repo,
	}
}

func (s *CatalogService) ListProducts(offset, limit int, category string, priceLessThan *float64) (*PaginatedResponse, error) {
	products, total, err := s.repo.GetProductsWithPagination(offset, limit, category, priceLessThan)
	if err != nil {
		return nil, err
	}

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


func (s *CatalogService) GetProductDetails(code string) (*ProductDetail, error) {
	product, err := s.repo.GetProductByCode(code)
	if err != nil {
		return nil, err
	}

	return mapProductToDetailDTO(product), nil
}
