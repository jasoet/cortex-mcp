package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

func setupFilmTestDB(t *testing.T) *gorm.DB {
	// Configure logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Silent,
			Colorful:      false,
		},
	)

	// Connect to test database
	// Using MySQL for testing - adjust connection string as needed
	dsn := "root:password@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		t.Skipf("Skipping test: Failed to connect to test database: %v", err)
		return nil
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&entity.Category{}, &entity.Film{}, &entity.Actor{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Clean up any existing data
	db.Exec("DELETE FROM film_actors")
	db.Exec("DELETE FROM film")
	db.Exec("DELETE FROM category")
	db.Exec("DELETE FROM actor")

	return db
}

func createTestCategory(t *testing.T, db *gorm.DB, name string) *entity.Category {
	category := &entity.Category{
		Name: name,
	}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category
}

func createTestFilm(t *testing.T, db *gorm.DB, title string, releaseYear int16, length int16, categoryID uint) *entity.Film {
	film := &entity.Film{
		Title:       title,
		ReleaseYear: releaseYear,
		Length:      length,
		CategoryID:  categoryID,
	}
	if err := db.Create(film).Error; err != nil {
		t.Fatalf("Failed to create test film: %v", err)
	}
	return film
}

func TestFilmRepository_Create(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Action")

	// Create a test film
	film := &entity.Film{
		Title:       "Test Film",
		ReleaseYear: 2023,
		Length:      120,
		CategoryID:  category.CategoryID,
	}

	// Test Create method
	err := repo.Create(ctx, film)
	if err != nil {
		t.Errorf("Failed to create film: %v", err)
	}
	if film.FilmID == 0 {
		t.Error("Expected FilmID to be non-zero after creation")
	}
}

func TestFilmRepository_FindByID(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Comedy")

	// Create a test film
	film := createTestFilm(t, db, "Test Comedy", 2022, 110, category.CategoryID)

	// Test FindByID method
	found, err := repo.FindByID(ctx, film.FilmID)
	if err != nil {
		t.Errorf("Failed to find film by ID: %v", err)
		return
	}
	if found == nil {
		t.Error("Expected to find film, but got nil")
		return
	}
	if found.Title != film.Title {
		t.Errorf("Expected film title %s, but got %s", film.Title, found.Title)
	}
}

func TestFilmRepository_FindAll(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Various")

	// Create test films
	films := []*entity.Film{
		createTestFilm(t, db, "Film 1", 2020, 100, category.CategoryID),
		createTestFilm(t, db, "Film 2", 2021, 110, category.CategoryID),
		createTestFilm(t, db, "Film 3", 2022, 120, category.CategoryID),
	}

	// Test FindAll method
	found, err := repo.FindAll(ctx)
	if err != nil {
		t.Errorf("Failed to find all films: %v", err)
		return
	}
	if len(found) != len(films) {
		t.Errorf("Expected to find %d films, but got %d", len(films), len(found))
	}
}

func TestFilmRepository_Update(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Adventure")

	// Create a test film
	film := createTestFilm(t, db, "Original Title", 2020, 100, category.CategoryID)

	// Update the film
	film.Title = "Updated Title"
	err := repo.Update(ctx, film)
	if err != nil {
		t.Errorf("Failed to update film: %v", err)
		return
	}

	// Verify the update
	found, err := repo.FindByID(ctx, film.FilmID)
	if err != nil {
		t.Errorf("Failed to find film by ID: %v", err)
		return
	}
	if found.Title != "Updated Title" {
		t.Errorf("Expected film title 'Updated Title', but got '%s'", found.Title)
	}
}

func TestFilmRepository_Delete(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Documentary")

	// Create a test film
	film := createTestFilm(t, db, "Film to Delete", 2021, 90, category.CategoryID)

	// Delete the film
	err := repo.Delete(ctx, film)
	if err != nil {
		t.Errorf("Failed to delete film: %v", err)
		return
	}

	// Verify the deletion
	_, err = repo.FindByID(ctx, film.FilmID)
	if err == nil {
		t.Error("Expected error when finding deleted film, but got nil")
	}
}

func TestFilmRepository_DeleteByID(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Animation")

	// Create a test film
	film := createTestFilm(t, db, "Film to Delete by ID", 2022, 95, category.CategoryID)

	// Delete the film by ID
	err := repo.DeleteByID(ctx, film.FilmID)
	if err != nil {
		t.Errorf("Failed to delete film by ID: %v", err)
		return
	}

	// Verify the deletion
	_, err = repo.FindByID(ctx, film.FilmID)
	if err == nil {
		t.Error("Expected error when finding deleted film, but got nil")
	}
}

func TestFilmRepository_FindByTitle(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Mixed")

	// Create test films
	createTestFilm(t, db, "Action Movie", 2020, 100, category.CategoryID)
	createTestFilm(t, db, "Action Adventure", 2021, 110, category.CategoryID)
	createTestFilm(t, db, "Comedy", 2022, 120, category.CategoryID)

	// Test FindByTitle method
	found, err := repo.FindByTitle(ctx, "Action")
	if err != nil {
		t.Errorf("Failed to find films by title: %v", err)
		return
	}
	if len(found) != 2 {
		t.Errorf("Expected to find 2 films, but got %d", len(found))
	}
}

func TestFilmRepository_FindByCategory(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create test categories
	actionCategory := createTestCategory(t, db, "Action")
	comedyCategory := createTestCategory(t, db, "Comedy")

	// Create test films
	createTestFilm(t, db, "Action Movie 1", 2020, 100, actionCategory.CategoryID)
	createTestFilm(t, db, "Action Movie 2", 2021, 110, actionCategory.CategoryID)
	createTestFilm(t, db, "Comedy Movie", 2022, 120, comedyCategory.CategoryID)

	// Test FindByCategory method
	found, err := repo.FindByCategory(ctx, actionCategory.CategoryID)
	if err != nil {
		t.Errorf("Failed to find films by category: %v", err)
		return
	}
	if len(found) != 2 {
		t.Errorf("Expected to find 2 films, but got %d", len(found))
	}
}

func TestFilmRepository_FindByReleaseYear(t *testing.T) {
	db := setupFilmTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewFilmRepository(db)
	ctx := context.Background()

	// Create a test category
	category := createTestCategory(t, db, "Various")

	// Create test films
	createTestFilm(t, db, "Film 2020", 2020, 100, category.CategoryID)
	createTestFilm(t, db, "Film 2021", 2021, 110, category.CategoryID)
	createTestFilm(t, db, "Another Film 2021", 2021, 120, category.CategoryID)

	// Test FindByReleaseYear method
	found, err := repo.FindByReleaseYear(ctx, 2021)
	if err != nil {
		t.Errorf("Failed to find films by release year: %v", err)
		return
	}
	if len(found) != 2 {
		t.Errorf("Expected to find 2 films, but got %d", len(found))
	}
}
