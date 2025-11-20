package repository

import (
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/internal/model"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() (model.Products, error) {
	products := make([]*model.Product, 0)
	if err := r.db.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
