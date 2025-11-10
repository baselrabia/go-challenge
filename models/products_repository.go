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

func (r *ProductsRepository) GetProductsWithPagination(offset, limit int) ([]Product, int64, error) {
	var products []Product
	var total int64

	// Get total count
	if err := r.db.Model(&Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated products
	if err := r.db.Offset(offset).Limit(limit).Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
