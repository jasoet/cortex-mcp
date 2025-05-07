package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// InventoryRepository is an interface for inventory operations
type InventoryRepository interface {
	Repository[entity.Inventory]

	// FindByFilm finds inventory items by film ID
	FindByFilm(ctx context.Context, filmID uint) ([]entity.Inventory, error)

	// FindByStore finds inventory items by store ID
	FindByStore(ctx context.Context, storeID uint) ([]entity.Inventory, error)

	// FindByFilmAndStore finds inventory items by film ID and store ID
	FindByFilmAndStore(ctx context.Context, filmID, storeID uint) ([]entity.Inventory, error)

	// FindAvailable finds inventory items that are not currently rented
	FindAvailable(ctx context.Context) ([]entity.Inventory, error)

	// FindAvailableByFilm finds available inventory items by film ID
	FindAvailableByFilm(ctx context.Context, filmID uint) ([]entity.Inventory, error)

	// FindAvailableByStore finds available inventory items by store ID
	FindAvailableByStore(ctx context.Context, storeID uint) ([]entity.Inventory, error)
}

// InventoryRepositoryImpl is an implementation of InventoryRepository
type InventoryRepositoryImpl struct {
	BaseRepository[entity.Inventory]
}

// NewInventoryRepository creates a new InventoryRepository
func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &InventoryRepositoryImpl{
		BaseRepository: BaseRepository[entity.Inventory]{
			DB: db,
		},
	}
}

// FindByFilm finds inventory items by film ID
func (r *InventoryRepositoryImpl) FindByFilm(ctx context.Context, filmID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).Where("film_id = ?", filmID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindByStore finds inventory items by store ID
func (r *InventoryRepositoryImpl) FindByStore(ctx context.Context, storeID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).Where("store_id = ?", storeID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindByFilmAndStore finds inventory items by film ID and store ID
func (r *InventoryRepositoryImpl) FindByFilmAndStore(ctx context.Context, filmID, storeID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).Where("film_id = ? AND store_id = ?", filmID, storeID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindAvailable finds inventory items that are not currently rented
func (r *InventoryRepositoryImpl) FindAvailable(ctx context.Context) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).
		Joins("LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL").
		Where("rental.rental_id IS NULL").
		Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindAvailableByFilm finds available inventory items by film ID
func (r *InventoryRepositoryImpl) FindAvailableByFilm(ctx context.Context, filmID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).
		Joins("LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL").
		Where("rental.rental_id IS NULL AND inventory.film_id = ?", filmID).
		Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

// FindAvailableByStore finds available inventory items by store ID
func (r *InventoryRepositoryImpl) FindAvailableByStore(ctx context.Context, storeID uint) ([]entity.Inventory, error) {
	var inventory []entity.Inventory
	if err := r.DB.WithContext(ctx).
		Joins("LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL").
		Where("rental.rental_id IS NULL AND inventory.store_id = ?", storeID).
		Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}
