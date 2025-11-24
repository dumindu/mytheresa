package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/internal/model"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetOneByCode(code string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Variants").Preload("Category").
		Where("code = ?", code).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAllWithFilter(limit, offset int64, filter *model.ProductFilter) (model.Products, error) {
	products := make([]*model.Product, 0)

	q := r.db.
		Model(&model.Product{}).
		Preload("Variants").
		Preload("Category").
		Joins(`
			LEFT JOIN categories
				ON categories.id = products.category_id
		`).
		Joins(`
			LEFT JOIN product_variants
				ON product_variants.product_id = products.id
		`)

	if filter != nil {
		if filter.Category != "" {
			q = q.Where("categories.name = ?", filter.Category)
		}

		if !filter.PriceLessThan.IsZero() {
			q = q.Where(`
				COALESCE(
					NULLIF(product_variants.price, 0),
					products.price
				) < ?
			`, filter.PriceLessThan)
		}
	}

	if limit > 0 {
		q = q.Limit(int(limit))
	}
	if offset > 0 {
		q = q.Offset(int(offset))
	}

	if err := q.Distinct().Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) CountAllWithFilter(filter *model.ProductFilter) (int64, error) {
	var count int64

	q := r.db.
		Model(&model.Product{}).
		Joins(`
			LEFT JOIN categories
				ON categories.id = products.category_id
		`).
		Joins(`
			LEFT JOIN product_variants
				ON product_variants.product_id = products.id
		`)

	if filter != nil {
		if filter.Category != "" {
			q = q.Where("categories.name = ?", filter.Category)
		}

		if !filter.PriceLessThan.IsZero() {
			q = q.Where(`
				COALESCE(
					NULLIF(product_variants.price, 0),
					products.price
				) < ?
			`, filter.PriceLessThan)
		}
	}

	if err := q.Distinct("products.id").Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
