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

func setupTestDB(t *testing.T) *gorm.DB {
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
	err = db.AutoMigrate(&entity.Category{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Clean up any existing data
	db.Exec("DELETE FROM category")

	return db
}

func TestCategoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create a test category
	category := &entity.Category{
		Name: "Action",
	}

	// Test Create method
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Failed to create category: %v", err)
	}
	if category.CategoryID == 0 {
		t.Error("Expected CategoryID to be non-zero after creation")
	}
}

func TestCategoryRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create a test category
	category := &entity.Category{
		Name: "Comedy",
	}
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Failed to create category: %v", err)
		return
	}

	// Test FindByID method
	found, err := repo.FindByID(ctx, category.CategoryID)
	if err != nil {
		t.Errorf("Failed to find category by ID: %v", err)
		return
	}
	if found == nil {
		t.Error("Expected to find category, but got nil")
		return
	}
	if found.Name != category.Name {
		t.Errorf("Expected category name %s, but got %s", category.Name, found.Name)
	}
}

func TestCategoryRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create test categories
	categories := []*entity.Category{
		{Name: "Drama"},
		{Name: "Horror"},
		{Name: "Sci-Fi"},
	}

	for _, cat := range categories {
		err := repo.Create(ctx, cat)
		if err != nil {
			t.Errorf("Failed to create category: %v", err)
			return
		}
	}

	// Test FindAll method
	found, err := repo.FindAll(ctx)
	if err != nil {
		t.Errorf("Failed to find all categories: %v", err)
		return
	}
	if len(found) != len(categories) {
		t.Errorf("Expected to find %d categories, but got %d", len(categories), len(found))
	}
}

func TestCategoryRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create a test category
	category := &entity.Category{
		Name: "Adventure",
	}
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Failed to create category: %v", err)
		return
	}

	// Update the category
	category.Name = "Updated Adventure"
	err = repo.Update(ctx, category)
	if err != nil {
		t.Errorf("Failed to update category: %v", err)
		return
	}

	// Verify the update
	found, err := repo.FindByID(ctx, category.CategoryID)
	if err != nil {
		t.Errorf("Failed to find category by ID: %v", err)
		return
	}
	if found.Name != "Updated Adventure" {
		t.Errorf("Expected category name 'Updated Adventure', but got '%s'", found.Name)
	}
}

func TestCategoryRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create a test category
	category := &entity.Category{
		Name: "Documentary",
	}
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Failed to create category: %v", err)
		return
	}

	// Delete the category
	err = repo.Delete(ctx, category)
	if err != nil {
		t.Errorf("Failed to delete category: %v", err)
		return
	}

	// Verify the deletion
	_, err = repo.FindByID(ctx, category.CategoryID)
	if err == nil {
		t.Error("Expected error when finding deleted category, but got nil")
	}
}

func TestCategoryRepository_DeleteByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create a test category
	category := &entity.Category{
		Name: "Animation",
	}
	err := repo.Create(ctx, category)
	if err != nil {
		t.Errorf("Failed to create category: %v", err)
		return
	}

	// Delete the category by ID
	err = repo.DeleteByID(ctx, category.CategoryID)
	if err != nil {
		t.Errorf("Failed to delete category by ID: %v", err)
		return
	}

	// Verify the deletion
	_, err = repo.FindByID(ctx, category.CategoryID)
	if err == nil {
		t.Error("Expected error when finding deleted category, but got nil")
	}
}

func TestCategoryRepository_FindByName(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return // Skip test if database connection failed
	}

	repo := NewCategoryRepository(db)
	ctx := context.Background()

	// Create test categories
	categories := []*entity.Category{
		{Name: "Action"},
		{Name: "Action Adventure"},
		{Name: "Comedy"},
	}

	for _, cat := range categories {
		err := repo.Create(ctx, cat)
		if err != nil {
			t.Errorf("Failed to create category: %v", err)
			return
		}
	}

	// Test FindByName method
	found, err := repo.FindByName(ctx, "Action")
	if err != nil {
		t.Errorf("Failed to find categories by name: %v", err)
		return
	}
	if len(found) != 2 {
		t.Errorf("Expected to find 2 categories, but got %d", len(found))
	}
}
