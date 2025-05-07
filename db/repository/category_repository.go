package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// CategoryRepository is an interface for category operations
type CategoryRepository interface {
	Repository[entity.Category]

	// FindByName finds categories by name
	FindByName(ctx context.Context, name string) ([]entity.Category, error)
}

// CategoryRepositoryImpl is an implementation of CategoryRepository
type CategoryRepositoryImpl struct {
	BaseRepository[entity.Category]
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &CategoryRepositoryImpl{
		BaseRepository: BaseRepository[entity.Category]{
			DB: db,
		},
	}
}

// FindByName finds categories by name
func (r *CategoryRepositoryImpl) FindByName(ctx context.Context, name string) ([]entity.Category, error) {
	var categories []entity.Category
	if err := r.DB.WithContext(ctx).Where("name LIKE ?", "%"+name+"%").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
