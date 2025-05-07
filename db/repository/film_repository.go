package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// FilmRepository is an interface for film operations
type FilmRepository interface {
	Repository[entity.Film]

	// FindByTitle finds films by title
	FindByTitle(ctx context.Context, title string) ([]entity.Film, error)

	// FindByCategory finds films by category ID
	FindByCategory(ctx context.Context, categoryID uint) ([]entity.Film, error)

	// FindByActor finds films by actor ID
	FindByActor(ctx context.Context, actorID uint) ([]entity.Film, error)

	// FindByReleaseYear finds films by release year
	FindByReleaseYear(ctx context.Context, year int16) ([]entity.Film, error)
}

// FilmRepositoryImpl is an implementation of FilmRepository
type FilmRepositoryImpl struct {
	BaseRepository[entity.Film]
}

// NewFilmRepository creates a new FilmRepository
func NewFilmRepository(db *gorm.DB) FilmRepository {
	return &FilmRepositoryImpl{
		BaseRepository: BaseRepository[entity.Film]{
			DB: db,
		},
	}
}

// FindByTitle finds films by title
func (r *FilmRepositoryImpl) FindByTitle(ctx context.Context, title string) ([]entity.Film, error) {
	var films []entity.Film
	if err := r.DB.WithContext(ctx).Where("title LIKE ?", "%"+title+"%").Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

// FindByCategory finds films by category ID
func (r *FilmRepositoryImpl) FindByCategory(ctx context.Context, categoryID uint) ([]entity.Film, error) {
	var films []entity.Film
	if err := r.DB.WithContext(ctx).Where("category_id = ?", categoryID).Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

// FindByActor finds films by actor ID
func (r *FilmRepositoryImpl) FindByActor(ctx context.Context, actorID uint) ([]entity.Film, error) {
	var films []entity.Film
	if err := r.DB.WithContext(ctx).Joins("JOIN film_actors ON film_actors.film_id = film.film_id").
		Where("film_actors.actor_id = ?", actorID).Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}

// FindByReleaseYear finds films by release year
func (r *FilmRepositoryImpl) FindByReleaseYear(ctx context.Context, year int16) ([]entity.Film, error) {
	var films []entity.Film
	if err := r.DB.WithContext(ctx).Where("release_year = ?", year).Find(&films).Error; err != nil {
		return nil, err
	}
	return films, nil
}
