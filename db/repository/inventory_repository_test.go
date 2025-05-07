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

func setupInventoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, InventoryRepository, func()) {
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

	repo := NewInventoryRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestInventoryRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	inventory := &entity.Inventory{
		FilmID:  1,
		StoreID: 1,
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `inventory`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			inventory.InventoryID,
			inventory.FilmID,
			inventory.StoreID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), inventory)
	if err != nil {
		t.Errorf("Error creating inventory: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventory
	expectedInventory := entity.Inventory{
		InventoryID: 1,
		FilmID:      1,
		StoreID:     1,
	}
	expectedInventory.ID = 1
	expectedInventory.CreatedAt = time.Now()
	expectedInventory.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	}).
		AddRow(
			expectedInventory.ID, expectedInventory.CreatedAt, expectedInventory.UpdatedAt, nil,
			expectedInventory.InventoryID, expectedInventory.FilmID, expectedInventory.StoreID,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` WHERE `inventory`.`id` = ? AND `inventory`.`deleted_at` IS NULL ORDER BY `inventory`.`id` LIMIT 1")).
		WithArgs(1).
		WillReturnRows(rows)

	inventory, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding inventory: %v", err)
	}

	if inventory.InventoryID != expectedInventory.InventoryID {
		t.Errorf("Expected InventoryID %d, got %d", expectedInventory.InventoryID, inventory.InventoryID)
	}

	if inventory.FilmID != expectedInventory.FilmID {
		t.Errorf("Expected FilmID %d, got %d", expectedInventory.FilmID, inventory.FilmID)
	}

	if inventory.StoreID != expectedInventory.StoreID {
		t.Errorf("Expected StoreID %d, got %d", expectedInventory.StoreID, inventory.StoreID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
		{
			InventoryID: 2,
			FilmID:      1,
			StoreID:     2,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()
	expectedInventories[1].ID = 2
	expectedInventories[1].CreatedAt = time.Now()
	expectedInventories[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` WHERE `inventory`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	inventories, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding inventories: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for i, inventory := range inventories {
		if inventory.InventoryID != expectedInventories[i].InventoryID {
			t.Errorf("Expected InventoryID %d, got %d", expectedInventories[i].InventoryID, inventory.InventoryID)
		}
		if inventory.FilmID != expectedInventories[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedInventories[i].FilmID, inventory.FilmID)
		}
		if inventory.StoreID != expectedInventories[i].StoreID {
			t.Errorf("Expected StoreID %d, got %d", expectedInventories[i].StoreID, inventory.StoreID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	inventory := &entity.Inventory{
		InventoryID: 1,
		FilmID:      1,
		StoreID:     1,
	}
	inventory.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			inventory.InventoryID,
			inventory.FilmID,
			inventory.StoreID,
			inventory.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), inventory)
	if err != nil {
		t.Errorf("Error updating inventory: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	inventory := &entity.Inventory{
		InventoryID: 1,
		FilmID:      1,
		StoreID:     1,
	}
	inventory.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `inventory` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			inventory.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), inventory)
	if err != nil {
		t.Errorf("Error deleting inventory: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `inventory` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			1,                // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error deleting inventory by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindByFilm(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
		{
			InventoryID: 2,
			FilmID:      1,
			StoreID:     2,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()
	expectedInventories[1].ID = 2
	expectedInventories[1].CreatedAt = time.Now()
	expectedInventories[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` WHERE film_id = ? AND `inventory`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	inventories, err := repo.FindByFilm(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding inventories by film: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for i, inventory := range inventories {
		if inventory.FilmID != expectedInventories[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedInventories[i].FilmID, inventory.FilmID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindByStore(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
		{
			InventoryID: 3,
			FilmID:      2,
			StoreID:     1,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()
	expectedInventories[1].ID = 3
	expectedInventories[1].CreatedAt = time.Now()
	expectedInventories[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` WHERE store_id = ? AND `inventory`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	inventories, err := repo.FindByStore(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding inventories by store: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for i, inventory := range inventories {
		if inventory.StoreID != expectedInventories[i].StoreID {
			t.Errorf("Expected StoreID %d, got %d", expectedInventories[i].StoreID, inventory.StoreID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindByFilmAndStore(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` WHERE film_id = ? AND store_id = ? AND `inventory`.`deleted_at` IS NULL")).
		WithArgs(uint(1), uint(1)).
		WillReturnRows(rows)

	inventories, err := repo.FindByFilmAndStore(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Error finding inventories by film and store: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for i, inventory := range inventories {
		if inventory.FilmID != expectedInventories[i].FilmID {
			t.Errorf("Expected FilmID %d, got %d", expectedInventories[i].FilmID, inventory.FilmID)
		}
		if inventory.StoreID != expectedInventories[i].StoreID {
			t.Errorf("Expected StoreID %d, got %d", expectedInventories[i].StoreID, inventory.StoreID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindAvailable(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
		{
			InventoryID: 2,
			FilmID:      1,
			StoreID:     2,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()
	expectedInventories[1].ID = 2
	expectedInventories[1].CreatedAt = time.Now()
	expectedInventories[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL WHERE rental.rental_id IS NULL AND `inventory`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	inventories, err := repo.FindAvailable(context.Background())
	if err != nil {
		t.Errorf("Error finding available inventories: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindAvailableByFilm(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL WHERE rental.rental_id IS NULL AND inventory.film_id = ? AND `inventory`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	inventories, err := repo.FindAvailableByFilm(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding available inventories by film: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for _, inventory := range inventories {
		if inventory.FilmID != 1 {
			t.Errorf("Expected FilmID 1, got %d", inventory.FilmID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestInventoryRepository_FindAvailableByStore(t *testing.T) {
	_, mock, repo, cleanup := setupInventoryTest(t)
	defer cleanup()

	// Define expected inventories
	expectedInventories := []entity.Inventory{
		{
			InventoryID: 1,
			FilmID:      1,
			StoreID:     1,
		},
	}
	expectedInventories[0].ID = 1
	expectedInventories[0].CreatedAt = time.Now()
	expectedInventories[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"inventory_id", "film_id", "store_id",
	})
	for _, inventory := range expectedInventories {
		rows.AddRow(
			inventory.ID, inventory.CreatedAt, inventory.UpdatedAt, nil,
			inventory.InventoryID, inventory.FilmID, inventory.StoreID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `inventory` LEFT JOIN rental ON rental.inventory_id = inventory.inventory_id AND rental.return_date IS NULL WHERE rental.rental_id IS NULL AND inventory.store_id = ? AND `inventory`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	inventories, err := repo.FindAvailableByStore(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding available inventories by store: %v", err)
	}

	if len(inventories) != len(expectedInventories) {
		t.Errorf("Expected %d inventories, got %d", len(expectedInventories), len(inventories))
	}

	for _, inventory := range inventories {
		if inventory.StoreID != 1 {
			t.Errorf("Expected StoreID 1, got %d", inventory.StoreID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
