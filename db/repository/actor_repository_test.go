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

func setupActorTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, ActorRepository, func()) {
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

	repo := NewActorRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestActorRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	actor := &entity.Actor{
		FirstName: "John",
		LastName:  "Doe",
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `actor`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			actor.FirstName,  // FirstName
			actor.LastName,   // LastName
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), actor)
	if err != nil {
		t.Errorf("Error creating actor: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	// Define expected actor
	expectedActor := entity.Actor{
		ActorID:   1,
		FirstName: "John",
		LastName:  "Doe",
	}
	expectedActor.ID = 1
	expectedActor.CreatedAt = time.Now()
	expectedActor.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "actor_id", "first_name", "last_name"}).
		AddRow(expectedActor.ID, expectedActor.CreatedAt, expectedActor.UpdatedAt, nil, expectedActor.ActorID, expectedActor.FirstName, expectedActor.LastName)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `actor` WHERE `actor`.`id` = ? AND `actor`.`deleted_at` IS NULL ORDER BY `actor`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	actor, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding actor: %v", err)
	}

	if actor.ActorID != expectedActor.ActorID {
		t.Errorf("Expected ActorID %d, got %d", expectedActor.ActorID, actor.ActorID)
	}

	if actor.FirstName != expectedActor.FirstName {
		t.Errorf("Expected FirstName %s, got %s", expectedActor.FirstName, actor.FirstName)
	}

	if actor.LastName != expectedActor.LastName {
		t.Errorf("Expected LastName %s, got %s", expectedActor.LastName, actor.LastName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	// Define expected actors
	expectedActors := []entity.Actor{
		{
			ActorID:   1,
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			ActorID:   2,
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}
	expectedActors[0].ID = 1
	expectedActors[0].CreatedAt = time.Now()
	expectedActors[0].UpdatedAt = time.Now()
	expectedActors[1].ID = 2
	expectedActors[1].CreatedAt = time.Now()
	expectedActors[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "actor_id", "first_name", "last_name"})
	for _, actor := range expectedActors {
		rows.AddRow(actor.ID, actor.CreatedAt, actor.UpdatedAt, nil, actor.ActorID, actor.FirstName, actor.LastName)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `actor` WHERE `actor`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	actors, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding actors: %v", err)
	}

	if len(actors) != len(expectedActors) {
		t.Errorf("Expected %d actors, got %d", len(expectedActors), len(actors))
	}

	for i, actor := range actors {
		if actor.ActorID != expectedActors[i].ActorID {
			t.Errorf("Expected ActorID %d, got %d", expectedActors[i].ActorID, actor.ActorID)
		}
		if actor.FirstName != expectedActors[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", expectedActors[i].FirstName, actor.FirstName)
		}
		if actor.LastName != expectedActors[i].LastName {
			t.Errorf("Expected LastName %s, got %s", expectedActors[i].LastName, actor.LastName)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	actor := &entity.Actor{
		ActorID:   1,
		FirstName: "John",
		LastName:  "Doe",
	}
	actor.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			actor.FirstName,
			actor.LastName,
			actor.ID,
			actor.ActorID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), actor)
	if err != nil {
		t.Errorf("Error updating actor: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	actor := &entity.Actor{
		ActorID:   1,
		FirstName: "John",
		LastName:  "Doe",
	}
	actor.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			actor.ID,
			actor.ActorID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), actor)
	if err != nil {
		t.Errorf("Error deleting actor: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
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
		t.Errorf("Error deleting actor by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_FindByName(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	// Define expected actors
	expectedActors := []entity.Actor{
		{
			ActorID:   1,
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	expectedActors[0].ID = 1
	expectedActors[0].CreatedAt = time.Now()
	expectedActors[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "actor_id", "first_name", "last_name"})
	for _, actor := range expectedActors {
		rows.AddRow(actor.ID, actor.CreatedAt, actor.UpdatedAt, nil, actor.ActorID, actor.FirstName, actor.LastName)
	}

	mock.ExpectQuery("SELECT").
		WithArgs("%John%", "%John%").
		WillReturnRows(rows)

	actors, err := repo.FindByName(context.Background(), "John")
	if err != nil {
		t.Errorf("Error finding actors by name: %v", err)
	}

	if len(actors) != len(expectedActors) {
		t.Errorf("Expected %d actors, got %d", len(expectedActors), len(actors))
	}

	for i, actor := range actors {
		if actor.ActorID != expectedActors[i].ActorID {
			t.Errorf("Expected ActorID %d, got %d", expectedActors[i].ActorID, actor.ActorID)
		}
		if actor.FirstName != expectedActors[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", expectedActors[i].FirstName, actor.FirstName)
		}
		if actor.LastName != expectedActors[i].LastName {
			t.Errorf("Expected LastName %s, got %s", expectedActors[i].LastName, actor.LastName)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_FindByFilm(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	// Define expected actors
	expectedActors := []entity.Actor{
		{
			ActorID:   1,
			FirstName: "John",
			LastName:  "Doe",
		},
		{
			ActorID:   2,
			FirstName: "Jane",
			LastName:  "Smith",
		},
	}
	expectedActors[0].ID = 1
	expectedActors[0].CreatedAt = time.Now()
	expectedActors[0].UpdatedAt = time.Now()
	expectedActors[1].ID = 2
	expectedActors[1].CreatedAt = time.Now()
	expectedActors[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "actor_id", "first_name", "last_name"})
	for _, actor := range expectedActors {
		rows.AddRow(actor.ID, actor.CreatedAt, actor.UpdatedAt, nil, actor.ActorID, actor.FirstName, actor.LastName)
	}

	mock.ExpectQuery("SELECT").
		WithArgs(1).
		WillReturnRows(rows)

	actors, err := repo.FindByFilm(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding actors by film: %v", err)
	}

	if len(actors) != len(expectedActors) {
		t.Errorf("Expected %d actors, got %d", len(expectedActors), len(actors))
	}

	for i, actor := range actors {
		if actor.ActorID != expectedActors[i].ActorID {
			t.Errorf("Expected ActorID %d, got %d", expectedActors[i].ActorID, actor.ActorID)
		}
		if actor.FirstName != expectedActors[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", expectedActors[i].FirstName, actor.FirstName)
		}
		if actor.LastName != expectedActors[i].LastName {
			t.Errorf("Expected LastName %s, got %s", expectedActors[i].LastName, actor.LastName)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestActorRepository_FindByID_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupActorTest(t)
	defer cleanup()

	// Expect the SELECT query with no results
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `actor` WHERE `actor`.`id` = ? AND `actor`.`deleted_at` IS NULL ORDER BY `actor`.`id` LIMIT ?")).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Error("Expected error when actor not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
