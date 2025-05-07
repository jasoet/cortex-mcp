package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// StoreRepository is an interface for store operations
type StoreRepository interface {
	Repository[entity.Store]

	// FindByName finds stores by name
	FindByName(ctx context.Context, name string) ([]entity.Store, error)

	// FindByCity finds stores by city
	FindByCity(ctx context.Context, city string) ([]entity.Store, error)

	// FindByCountry finds stores by country
	FindByCountry(ctx context.Context, country string) ([]entity.Store, error)
}

// StoreRepositoryImpl is an implementation of StoreRepository
type StoreRepositoryImpl struct {
	BaseRepository[entity.Store]
}

// NewStoreRepository creates a new StoreRepository
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &StoreRepositoryImpl{
		BaseRepository: BaseRepository[entity.Store]{
			DB: db,
		},
	}
}

// FindByName finds stores by name
func (r *StoreRepositoryImpl) FindByName(ctx context.Context, name string) ([]entity.Store, error) {
	var stores []entity.Store
	if err := r.DB.WithContext(ctx).Where("store_name LIKE ?", "%"+name+"%").Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

// FindByCity finds stores by city
func (r *StoreRepositoryImpl) FindByCity(ctx context.Context, city string) ([]entity.Store, error) {
	var stores []entity.Store
	if err := r.DB.WithContext(ctx).Where("city LIKE ?", "%"+city+"%").Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}

// FindByCountry finds stores by country
func (r *StoreRepositoryImpl) FindByCountry(ctx context.Context, country string) ([]entity.Store, error) {
	var stores []entity.Store
	if err := r.DB.WithContext(ctx).Where("country LIKE ?", "%"+country+"%").Find(&stores).Error; err != nil {
		return nil, err
	}
	return stores, nil
}
