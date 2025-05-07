package repository

import (
	"context"
	"gorm.io/gorm"
)

// Repository is a generic interface for database operations
type Repository[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// FindByID finds an entity by its ID
	FindByID(ctx context.Context, id uint) (*T, error)

	// FindAll returns all entities
	FindAll(ctx context.Context) ([]T, error)

	// Update updates an entity
	Update(ctx context.Context, entity *T) error

	// Delete deletes an entity
	Delete(ctx context.Context, entity *T) error

	// DeleteByID deletes an entity by its ID
	DeleteByID(ctx context.Context, id uint) error
}

// BaseRepository is a base implementation of the Repository interface
type BaseRepository[T any] struct {
	DB *gorm.DB
}

// Create creates a new entity
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// FindByID finds an entity by its ID
func (r *BaseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	var entity T
	if err := r.DB.WithContext(ctx).First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll returns all entities
func (r *BaseRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.DB.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Update updates an entity
func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete deletes an entity
func (r *BaseRepository[T]) Delete(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Delete(entity).Error
}

// DeleteByID deletes an entity by its ID
func (r *BaseRepository[T]) DeleteByID(ctx context.Context, id uint) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, id).Error
}
