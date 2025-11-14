package models

import (
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) GetProductsWithPagination(offset, limit int, category string, priceLessThan *float64) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := r.db.Model(&Product{})

	// category filter
	if category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.code = ?", category)
	}

	// price filter
	if priceLessThan != nil {
		query = query.Where("products.price < ?", *priceLessThan)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Where("code = ?", code).Preload("Category").Preload("Variants").First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
