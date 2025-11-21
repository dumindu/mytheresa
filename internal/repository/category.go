package repository

import (
	"gorm.io/gorm"

	"github.com/mytheresa/go-hiring-challenge/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(category *model.Category) (*model.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) GetAll(limit, offset int64) (model.Categories, error) {
	categories := make([]*model.Category, 0)
	q := r.db

	if limit > 0 {
		q = q.Limit(int(limit))
	}
	if offset > 0 {
		q = q.Offset(int(offset))
	}

	if err := q.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Category{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
