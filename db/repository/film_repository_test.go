package repository

import (
	"CortexMCP/db/entity"
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupFilmTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, FilmRepository, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	repo := NewFilmRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestFilmRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	film := &entity.Film{
		Title:       "The Matrix",
		ReleaseYear: 1999,
		Length:      136,
		CategoryID:  1,
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `film`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			film.Title,       // Title
			film.ReleaseYear, // ReleaseYear
			film.Length,      // Length
			film.CategoryID,  // CategoryID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), film)
	if err != nil {
		t.Errorf("Error creating film: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected film
	expectedFilm := entity.Film{
		FilmID:      1,
		Title:       "The Matrix",
		ReleaseYear: 1999,
		Length:      136,
		CategoryID:  1,
	}
	expectedFilm.ID = 1
	expectedFilm.CreatedAt = time.Now()
	expectedFilm.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"}).
		AddRow(expectedFilm.ID, expectedFilm.CreatedAt, expectedFilm.UpdatedAt, nil, expectedFilm.FilmID, expectedFilm.Title, expectedFilm.ReleaseYear, expectedFilm.Length, expectedFilm.CategoryID)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `film` WHERE `film`.`id` = ? AND `film`.`deleted_at` IS NULL ORDER BY `film`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	film, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding film: %v", err)
	}

	if film.FilmID != expectedFilm.FilmID {
		t.Errorf("Expected FilmID %d, got %d", expectedFilm.FilmID, film.FilmID)
	}

	if film.Title != expectedFilm.Title {
		t.Errorf("Expected Title %s, got %s", expectedFilm.Title, film.Title)
	}

	if film.ReleaseYear != expectedFilm.ReleaseYear {
		t.Errorf("Expected ReleaseYear %d, got %d", expectedFilm.ReleaseYear, film.ReleaseYear)
	}

	if film.Length != expectedFilm.Length {
		t.Errorf("Expected Length %d, got %d", expectedFilm.Length, film.Length)
	}

	if film.CategoryID != expectedFilm.CategoryID {
		t.Errorf("Expected CategoryID %d, got %d", expectedFilm.CategoryID, film.CategoryID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected films
	expectedFilms := []entity.Film{
		{
			FilmID:      1,
			Title:       "The Matrix",
			ReleaseYear: 1999,
			Length:      136,
			CategoryID:  1,
		},
		{
			FilmID:      2,
			Title:       "The Matrix Reloaded",
			ReleaseYear: 2003,
			Length:      138,
			CategoryID:  1,
		},
	}
	expectedFilms[0].ID = 1
	expectedFilms[0].CreatedAt = time.Now()
	expectedFilms[0].UpdatedAt = time.Now()
	expectedFilms[1].ID = 2
	expectedFilms[1].CreatedAt = time.Now()
	expectedFilms[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"})
	for _, film := range expectedFilms {
		rows.AddRow(film.ID, film.CreatedAt, film.UpdatedAt, nil, film.FilmID, film.Title, film.ReleaseYear, film.Length, film.CategoryID)
	}

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	films, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding films: %v", err)
	}

	if len(films) != len(expectedFilms) {
		t.Errorf("Expected %d films, got %d", len(expectedFilms), len(films))
	}

	for i, film := range films {
		if film.FilmID != expectedFilms[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedFilms[i].FilmID, film.FilmID)
		}
		if film.Title != expectedFilms[i].Title {
			t.Errorf("Expected Title %s, got %s", expectedFilms[i].Title, film.Title)
		}
		if film.ReleaseYear != expectedFilms[i].ReleaseYear {
			t.Errorf("Expected ReleaseYear %d, got %d", expectedFilms[i].ReleaseYear, film.ReleaseYear)
		}
		if film.Length != expectedFilms[i].Length {
			t.Errorf("Expected Length %d, got %d", expectedFilms[i].Length, film.Length)
		}
		if film.CategoryID != expectedFilms[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedFilms[i].CategoryID, film.CategoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	film := &entity.Film{
		FilmID:      1,
		Title:       "The Matrix",
		ReleaseYear: 1999,
		Length:      136,
		CategoryID:  1,
	}
	film.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			film.Title,       // Title
			film.ReleaseYear, // ReleaseYear
			film.Length,      // Length
			film.CategoryID,  // CategoryID
			film.ID,          // ID
			film.FilmID,      // FilmID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), film)
	if err != nil {
		t.Errorf("Error updating film: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	film := &entity.Film{
		FilmID:      1,
		Title:       "The Matrix",
		ReleaseYear: 1999,
		Length:      136,
		CategoryID:  1,
	}
	film.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			film.ID,          // ID
			film.FilmID,      // FilmID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), film)
	if err != nil {
		t.Errorf("Error deleting film: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			1,                // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error deleting film by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByTitle(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected films
	expectedFilms := []entity.Film{
		{
			FilmID:      1,
			Title:       "The Matrix",
			ReleaseYear: 1999,
			Length:      136,
			CategoryID:  1,
		},
	}
	expectedFilms[0].ID = 1
	expectedFilms[0].CreatedAt = time.Now()
	expectedFilms[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"})
	for _, film := range expectedFilms {
		rows.AddRow(film.ID, film.CreatedAt, film.UpdatedAt, nil, film.FilmID, film.Title, film.ReleaseYear, film.Length, film.CategoryID)
	}

	mock.ExpectQuery("SELECT").
		WithArgs("%Matrix%").
		WillReturnRows(rows)

	films, err := repo.FindByTitle(context.Background(), "Matrix")
	if err != nil {
		t.Errorf("Error finding films by title: %v", err)
	}

	if len(films) != len(expectedFilms) {
		t.Errorf("Expected %d films, got %d", len(expectedFilms), len(films))
	}

	for i, film := range films {
		if film.FilmID != expectedFilms[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedFilms[i].FilmID, film.FilmID)
		}
		if film.Title != expectedFilms[i].Title {
			t.Errorf("Expected Title %s, got %s", expectedFilms[i].Title, film.Title)
		}
		if film.ReleaseYear != expectedFilms[i].ReleaseYear {
			t.Errorf("Expected ReleaseYear %d, got %d", expectedFilms[i].ReleaseYear, film.ReleaseYear)
		}
		if film.Length != expectedFilms[i].Length {
			t.Errorf("Expected Length %d, got %d", expectedFilms[i].Length, film.Length)
		}
		if film.CategoryID != expectedFilms[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedFilms[i].CategoryID, film.CategoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByCategory(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected films
	expectedFilms := []entity.Film{
		{
			FilmID:      1,
			Title:       "The Matrix",
			ReleaseYear: 1999,
			Length:      136,
			CategoryID:  1,
		},
		{
			FilmID:      2,
			Title:       "The Matrix Reloaded",
			ReleaseYear: 2003,
			Length:      138,
			CategoryID:  1,
		},
	}
	expectedFilms[0].ID = 1
	expectedFilms[0].CreatedAt = time.Now()
	expectedFilms[0].UpdatedAt = time.Now()
	expectedFilms[1].ID = 2
	expectedFilms[1].CreatedAt = time.Now()
	expectedFilms[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"})
	for _, film := range expectedFilms {
		rows.AddRow(film.ID, film.CreatedAt, film.UpdatedAt, nil, film.FilmID, film.Title, film.ReleaseYear, film.Length, film.CategoryID)
	}

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	films, err := repo.FindByCategory(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding films by category: %v", err)
	}

	if len(films) != len(expectedFilms) {
		t.Errorf("Expected %d films, got %d", len(expectedFilms), len(films))
	}

	for i, film := range films {
		if film.FilmID != expectedFilms[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedFilms[i].FilmID, film.FilmID)
		}
		if film.Title != expectedFilms[i].Title {
			t.Errorf("Expected Title %s, got %s", expectedFilms[i].Title, film.Title)
		}
		if film.ReleaseYear != expectedFilms[i].ReleaseYear {
			t.Errorf("Expected ReleaseYear %d, got %d", expectedFilms[i].ReleaseYear, film.ReleaseYear)
		}
		if film.Length != expectedFilms[i].Length {
			t.Errorf("Expected Length %d, got %d", expectedFilms[i].Length, film.Length)
		}
		if film.CategoryID != expectedFilms[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedFilms[i].CategoryID, film.CategoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByActor(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected films
	expectedFilms := []entity.Film{
		{
			FilmID:      1,
			Title:       "The Matrix",
			ReleaseYear: 1999,
			Length:      136,
			CategoryID:  1,
		},
		{
			FilmID:      2,
			Title:       "The Matrix Reloaded",
			ReleaseYear: 2003,
			Length:      138,
			CategoryID:  1,
		},
	}
	expectedFilms[0].ID = 1
	expectedFilms[0].CreatedAt = time.Now()
	expectedFilms[0].UpdatedAt = time.Now()
	expectedFilms[1].ID = 2
	expectedFilms[1].CreatedAt = time.Now()
	expectedFilms[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"})
	for _, film := range expectedFilms {
		rows.AddRow(film.ID, film.CreatedAt, film.UpdatedAt, nil, film.FilmID, film.Title, film.ReleaseYear, film.Length, film.CategoryID)
	}

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	films, err := repo.FindByActor(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding films by actor: %v", err)
	}

	if len(films) != len(expectedFilms) {
		t.Errorf("Expected %d films, got %d", len(expectedFilms), len(films))
	}

	for i, film := range films {
		if film.FilmID != expectedFilms[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedFilms[i].FilmID, film.FilmID)
		}
		if film.Title != expectedFilms[i].Title {
			t.Errorf("Expected Title %s, got %s", expectedFilms[i].Title, film.Title)
		}
		if film.ReleaseYear != expectedFilms[i].ReleaseYear {
			t.Errorf("Expected ReleaseYear %d, got %d", expectedFilms[i].ReleaseYear, film.ReleaseYear)
		}
		if film.Length != expectedFilms[i].Length {
			t.Errorf("Expected Length %d, got %d", expectedFilms[i].Length, film.Length)
		}
		if film.CategoryID != expectedFilms[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedFilms[i].CategoryID, film.CategoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByReleaseYear(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Define expected films
	expectedFilms := []entity.Film{
		{
			FilmID:      1,
			Title:       "The Matrix",
			ReleaseYear: 1999,
			Length:      136,
			CategoryID:  1,
		},
	}
	expectedFilms[0].ID = 1
	expectedFilms[0].CreatedAt = time.Now()
	expectedFilms[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "film_id", "title", "release_year", "length", "category_id"})
	for _, film := range expectedFilms {
		rows.AddRow(film.ID, film.CreatedAt, film.UpdatedAt, nil, film.FilmID, film.Title, film.ReleaseYear, film.Length, film.CategoryID)
	}

	mock.ExpectQuery("SELECT").
		WithArgs(int16(1999)).
		WillReturnRows(rows)

	films, err := repo.FindByReleaseYear(context.Background(), 1999)
	if err != nil {
		t.Errorf("Error finding films by release year: %v", err)
	}

	if len(films) != len(expectedFilms) {
		t.Errorf("Expected %d films, got %d", len(expectedFilms), len(films))
	}

	for i, film := range films {
		if film.FilmID != expectedFilms[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedFilms[i].FilmID, film.FilmID)
		}
		if film.Title != expectedFilms[i].Title {
			t.Errorf("Expected Title %s, got %s", expectedFilms[i].Title, film.Title)
		}
		if film.ReleaseYear != expectedFilms[i].ReleaseYear {
			t.Errorf("Expected ReleaseYear %d, got %d", expectedFilms[i].ReleaseYear, film.ReleaseYear)
		}
		if film.Length != expectedFilms[i].Length {
			t.Errorf("Expected Length %d, got %d", expectedFilms[i].Length, film.Length)
		}
		if film.CategoryID != expectedFilms[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedFilms[i].CategoryID, film.CategoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestFilmRepository_FindByID_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupFilmTest(t)
	defer cleanup()

	// Expect the SELECT query with no results
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `film` WHERE `film`.`id` = ? AND `film`.`deleted_at` IS NULL ORDER BY `film`.`id` LIMIT ?")).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Error("Expected error when film not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
