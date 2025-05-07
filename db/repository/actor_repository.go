package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// ActorRepository is an interface for actor operations
type ActorRepository interface {
	Repository[entity.Actor]

	// FindByName finds actors by first or last name
	FindByName(ctx context.Context, name string) ([]entity.Actor, error)

	// FindByFilm finds actors by film ID
	FindByFilm(ctx context.Context, filmID uint) ([]entity.Actor, error)
}

// ActorRepositoryImpl is an implementation of ActorRepository
type ActorRepositoryImpl struct {
	BaseRepository[entity.Actor]
}

// NewActorRepository creates a new ActorRepository
func NewActorRepository(db *gorm.DB) ActorRepository {
	return &ActorRepositoryImpl{
		BaseRepository: BaseRepository[entity.Actor]{
			DB: db,
		},
	}
}

// FindByName finds actors by first or last name
func (r *ActorRepositoryImpl) FindByName(ctx context.Context, name string) ([]entity.Actor, error) {
	var actors []entity.Actor
	if err := r.DB.WithContext(ctx).Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%").Find(&actors).Error; err != nil {
		return nil, err
	}
	return actors, nil
}

// FindByFilm finds actors by film ID
func (r *ActorRepositoryImpl) FindByFilm(ctx context.Context, filmID uint) ([]entity.Actor, error) {
	var actors []entity.Actor
	if err := r.DB.WithContext(ctx).Joins("JOIN film_actors ON film_actors.actor_id = actor.actor_id").
		Where("film_actors.film_id = ?", filmID).Find(&actors).Error; err != nil {
		return nil, err
	}
	return actors, nil
}
