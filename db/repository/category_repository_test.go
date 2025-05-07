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

func setupCategoryTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, CategoryRepository, func()) {
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

	repo := NewCategoryRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	category := &entity.Category{
		Name: "Action",
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `category`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			category.Name,    // Name
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), category)
	if err != nil {
		t.Errorf("Error creating category: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	// Define expected category
	expectedCategory := entity.Category{
		CategoryID: 1,
		Name:       "Action",
	}
	expectedCategory.ID = 1
	expectedCategory.CreatedAt = time.Now()
	expectedCategory.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "category_id", "name"}).
		AddRow(expectedCategory.ID, expectedCategory.CreatedAt, expectedCategory.UpdatedAt, nil, expectedCategory.CategoryID, expectedCategory.Name)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `category` WHERE `category`.`id` = ? AND `category`.`deleted_at` IS NULL ORDER BY `category`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	category, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding category: %v", err)
	}

	if category.CategoryID != expectedCategory.CategoryID {
		t.Errorf("Expected CategoryID %d, got %d", expectedCategory.CategoryID, category.CategoryID)
	}

	if category.Name != expectedCategory.Name {
		t.Errorf("Expected Name %s, got %s", expectedCategory.Name, category.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	// Define expected categories
	expectedCategories := []entity.Category{
		{
			CategoryID: 1,
			Name:       "Action",
		},
		{
			CategoryID: 2,
			Name:       "Comedy",
		},
	}
	expectedCategories[0].ID = 1
	expectedCategories[0].CreatedAt = time.Now()
	expectedCategories[0].UpdatedAt = time.Now()
	expectedCategories[1].ID = 2
	expectedCategories[1].CreatedAt = time.Now()
	expectedCategories[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "category_id", "name"})
	for _, category := range expectedCategories {
		rows.AddRow(category.ID, category.CreatedAt, category.UpdatedAt, nil, category.CategoryID, category.Name)
	}

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	categories, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding categories: %v", err)
	}

	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}

	for i, category := range categories {
		if category.CategoryID != expectedCategories[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedCategories[i].CategoryID, category.CategoryID)
		}
		if category.Name != expectedCategories[i].Name {
			t.Errorf("Expected Name %s, got %s", expectedCategories[i].Name, category.Name)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	category := &entity.Category{
		CategoryID: 1,
		Name:       "Action",
	}
	category.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(),    // CreatedAt
			sqlmock.AnyArg(),    // UpdatedAt
			sqlmock.AnyArg(),    // DeletedAt
			category.Name,       // Name
			category.ID,         // ID
			category.CategoryID, // CategoryID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), category)
	if err != nil {
		t.Errorf("Error updating category: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	category := &entity.Category{
		CategoryID: 1,
		Name:       "Action",
	}
	category.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			category.ID,
			category.CategoryID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), category)
	if err != nil {
		t.Errorf("Error deleting category: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
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
		t.Errorf("Error deleting category by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByName(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	// Define expected categories
	expectedCategories := []entity.Category{
		{
			CategoryID: 1,
			Name:       "Action",
		},
	}
	expectedCategories[0].ID = 1
	expectedCategories[0].CreatedAt = time.Now()
	expectedCategories[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "category_id", "name"})
	for _, category := range expectedCategories {
		rows.AddRow(category.ID, category.CreatedAt, category.UpdatedAt, nil, category.CategoryID, category.Name)
	}

	mock.ExpectQuery("SELECT").
		WithArgs("%Action%").
		WillReturnRows(rows)

	categories, err := repo.FindByName(context.Background(), "Action")
	if err != nil {
		t.Errorf("Error finding categories by name: %v", err)
	}

	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}

	for i, category := range categories {
		if category.CategoryID != expectedCategories[i].CategoryID {
			t.Errorf("Expected CategoryID %d, got %d", expectedCategories[i].CategoryID, category.CategoryID)
		}
		if category.Name != expectedCategories[i].Name {
			t.Errorf("Expected Name %s, got %s", expectedCategories[i].Name, category.Name)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCategoryRepository_FindByID_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupCategoryTest(t)
	defer cleanup()

	// Expect the SELECT query with no results
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `category` WHERE `category`.`id` = ? AND `category`.`deleted_at` IS NULL ORDER BY `category`.`id` LIMIT ?")).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Error("Expected error when category not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
