package catalog

import "github.com/mytheresa/go-hiring-challenge/models"

func mapProductToDTO(p models.Product) Product {
	dto := Product{
		Code:  p.Code,
		Price: p.Price.InexactFloat64(),
	}

	if p.Category != nil {
		dto.Category = &Category{
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	return dto
}


func mapProductToDetailDTO(p *models.Product) *ProductDetail {
	productPrice := p.Price.InexactFloat64()

	detail := &ProductDetail{
		Code:  p.Code,
		Price: productPrice,
	}

	if p.Category != nil {
		detail.Category = &Category{
			Code: p.Category.Code,
			Name: p.Category.Name,
		}
	}

	// Map variants with price inheritance logic
	variants := make([]VariantDetail, len(p.Variants))
	for i, v := range p.Variants {
		variantPrice := productPrice

		// Use variant specific price if set (non-zero)
		if !v.Price.IsZero() {
			variantPrice = v.Price.InexactFloat64()
		}

		variants[i] = VariantDetail{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: variantPrice,
		}
	}
	detail.Variants = variants

	return detail
}
