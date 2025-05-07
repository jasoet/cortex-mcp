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

func setupRentalTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, RentalRepository, func()) {
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

	repo := NewRentalRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestRentalRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	rentalDate := time.Now()
	returnDate := rentalDate.Add(time.Hour * 24 * 3) // 3 days later
	rental := &entity.Rental{
		RentalDate:  rentalDate,
		InventoryID: 1,
		CustomerID:  1,
		ReturnDate:  &returnDate,
		StaffID:     1,
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `rental`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			rental.RentalID,
			rental.RentalDate,
			rental.InventoryID,
			rental.CustomerID,
			rental.ReturnDate,
			rental.StaffID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), rental)
	if err != nil {
		t.Errorf("Error creating rental: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rental
	rentalDate := time.Now()
	returnDate := rentalDate.Add(time.Hour * 24 * 3) // 3 days later
	expectedRental := entity.Rental{
		RentalID:    1,
		RentalDate:  rentalDate,
		InventoryID: 1,
		CustomerID:  1,
		ReturnDate:  &returnDate,
		StaffID:     1,
	}
	expectedRental.ID = 1
	expectedRental.CreatedAt = time.Now()
	expectedRental.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	}).
		AddRow(
			expectedRental.ID, expectedRental.CreatedAt, expectedRental.UpdatedAt, nil,
			expectedRental.RentalID, expectedRental.RentalDate, expectedRental.InventoryID,
			expectedRental.CustomerID, expectedRental.ReturnDate, expectedRental.StaffID,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE `rental`.`id` = ? AND `rental`.`deleted_at` IS NULL ORDER BY `rental`.`id` LIMIT 1")).
		WithArgs(1).
		WillReturnRows(rows)

	rental, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding rental: %v", err)
	}

	if rental.RentalID != expectedRental.RentalID {
		t.Errorf("Expected RentalID %d, got %d", expectedRental.RentalID, rental.RentalID)
	}

	if rental.InventoryID != expectedRental.InventoryID {
		t.Errorf("Expected InventoryID %d, got %d", expectedRental.InventoryID, rental.InventoryID)
	}

	if rental.CustomerID != expectedRental.CustomerID {
		t.Errorf("Expected CustomerID %d, got %d", expectedRental.CustomerID, rental.CustomerID)
	}

	if rental.StaffID != expectedRental.StaffID {
		t.Errorf("Expected StaffID %d, got %d", expectedRental.StaffID, rental.StaffID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now()
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3) // 3 days later
	rentalDate2 := time.Now().Add(time.Hour * 24)      // 1 day later
	returnDate2 := rentalDate2.Add(time.Hour * 24 * 5) // 5 days later
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    2,
			RentalDate:  rentalDate2,
			InventoryID: 2,
			CustomerID:  2,
			ReturnDate:  &returnDate2,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 2
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE `rental`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	rentals, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding rentals: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for i, rental := range rentals {
		if rental.RentalID != expectedRentals[i].RentalID {
			t.Errorf("Expected RentalID %d, got %d", expectedRentals[i].RentalID, rental.RentalID)
		}
		if rental.InventoryID != expectedRentals[i].InventoryID {
			t.Errorf("Expected InventoryID %d, got %d", expectedRentals[i].InventoryID, rental.InventoryID)
		}
		if rental.CustomerID != expectedRentals[i].CustomerID {
			t.Errorf("Expected CustomerID %d, got %d", expectedRentals[i].CustomerID, rental.CustomerID)
		}
		if rental.StaffID != expectedRentals[i].StaffID {
			t.Errorf("Expected StaffID %d, got %d", expectedRentals[i].StaffID, rental.StaffID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	rentalDate := time.Now()
	returnDate := rentalDate.Add(time.Hour * 24 * 3) // 3 days later
	rental := &entity.Rental{
		RentalID:    1,
		RentalDate:  rentalDate,
		InventoryID: 1,
		CustomerID:  1,
		ReturnDate:  &returnDate,
		StaffID:     1,
	}
	rental.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			rental.RentalID,
			rental.RentalDate,
			rental.InventoryID,
			rental.CustomerID,
			rental.ReturnDate,
			rental.StaffID,
			rental.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), rental)
	if err != nil {
		t.Errorf("Error updating rental: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	rental := &entity.Rental{
		RentalID:    1,
		InventoryID: 1,
		CustomerID:  1,
		StaffID:     1,
	}
	rental.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `rental` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			rental.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), rental)
	if err != nil {
		t.Errorf("Error deleting rental: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `rental` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			1,                // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error deleting rental by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindByCustomer(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now()
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3) // 3 days later
	rentalDate2 := time.Now().Add(time.Hour * 24)      // 1 day later
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    3,
			RentalDate:  rentalDate2,
			InventoryID: 3,
			CustomerID:  1,
			ReturnDate:  nil,
			StaffID:     2,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 3
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE customer_id = ? AND `rental`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	rentals, err := repo.FindByCustomer(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding rentals by customer: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.CustomerID != 1 {
			t.Errorf("Expected CustomerID 1, got %d", rental.CustomerID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindByStaff(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now()
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3) // 3 days later
	rentalDate2 := time.Now().Add(time.Hour * 24)      // 1 day later
	returnDate2 := rentalDate2.Add(time.Hour * 24 * 5) // 5 days later
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    2,
			RentalDate:  rentalDate2,
			InventoryID: 2,
			CustomerID:  2,
			ReturnDate:  &returnDate2,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 2
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE staff_id = ? AND `rental`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	rentals, err := repo.FindByStaff(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding rentals by staff: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.StaffID != 1 {
			t.Errorf("Expected StaffID 1, got %d", rental.StaffID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindByInventory(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now()
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3) // 3 days later
	rentalDate2 := time.Now().Add(time.Hour * 24 * 7)  // 7 days later
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    4,
			RentalDate:  rentalDate2,
			InventoryID: 1,
			CustomerID:  2,
			ReturnDate:  nil,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 4
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE inventory_id = ? AND `rental`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	rentals, err := repo.FindByInventory(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding rentals by inventory: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.InventoryID != 1 {
			t.Errorf("Expected InventoryID 1, got %d", rental.InventoryID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindByDateRange(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define date range
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)

	// Define expected rentals
	rentalDate1 := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3) // 3 days later
	rentalDate2 := time.Date(2023, 1, 20, 14, 30, 0, 0, time.UTC)
	returnDate2 := rentalDate2.Add(time.Hour * 24 * 5) // 5 days later
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    2,
			RentalDate:  rentalDate2,
			InventoryID: 2,
			CustomerID:  2,
			ReturnDate:  &returnDate2,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 2
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE rental_date BETWEEN ? AND ? AND `rental`.`deleted_at` IS NULL")).
		WithArgs(startDate, endDate).
		WillReturnRows(rows)

	rentals, err := repo.FindByDateRange(context.Background(), startDate, endDate)
	if err != nil {
		t.Errorf("Error finding rentals by date range: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindOverdue(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now().Add(-time.Hour * 24 * 10) // 10 days ago
	rentalDate2 := time.Now().Add(-time.Hour * 24 * 8)  // 8 days ago
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  nil,
			StaffID:     1,
		},
		{
			RentalID:    2,
			RentalDate:  rentalDate2,
			InventoryID: 2,
			CustomerID:  2,
			ReturnDate:  nil,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 2
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE return_date IS NULL AND rental_date < ? AND `rental`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	rentals, err := repo.FindOverdue(context.Background(), 7) // Overdue if rented more than 7 days ago
	if err != nil {
		t.Errorf("Error finding overdue rentals: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.ReturnDate != nil {
			t.Errorf("Expected ReturnDate to be nil, got %v", rental.ReturnDate)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindReturned(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now().Add(-time.Hour * 24 * 10) // 10 days ago
	returnDate1 := rentalDate1.Add(time.Hour * 24 * 3)  // 3 days after rental
	rentalDate2 := time.Now().Add(-time.Hour * 24 * 8)  // 8 days ago
	returnDate2 := rentalDate2.Add(time.Hour * 24 * 2)  // 2 days after rental
	expectedRentals := []entity.Rental{
		{
			RentalID:    1,
			RentalDate:  rentalDate1,
			InventoryID: 1,
			CustomerID:  1,
			ReturnDate:  &returnDate1,
			StaffID:     1,
		},
		{
			RentalID:    2,
			RentalDate:  rentalDate2,
			InventoryID: 2,
			CustomerID:  2,
			ReturnDate:  &returnDate2,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 1
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 2
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE return_date IS NOT NULL AND `rental`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	rentals, err := repo.FindReturned(context.Background())
	if err != nil {
		t.Errorf("Error finding returned rentals: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.ReturnDate == nil {
			t.Errorf("Expected ReturnDate to be non-nil, got nil")
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestRentalRepository_FindNotReturned(t *testing.T) {
	_, mock, repo, cleanup := setupRentalTest(t)
	defer cleanup()

	// Define expected rentals
	rentalDate1 := time.Now().Add(-time.Hour * 24 * 5) // 5 days ago
	rentalDate2 := time.Now().Add(-time.Hour * 24 * 3) // 3 days ago
	expectedRentals := []entity.Rental{
		{
			RentalID:    3,
			RentalDate:  rentalDate1,
			InventoryID: 3,
			CustomerID:  1,
			ReturnDate:  nil,
			StaffID:     2,
		},
		{
			RentalID:    4,
			RentalDate:  rentalDate2,
			InventoryID: 4,
			CustomerID:  2,
			ReturnDate:  nil,
			StaffID:     1,
		},
	}
	expectedRentals[0].ID = 3
	expectedRentals[0].CreatedAt = time.Now()
	expectedRentals[0].UpdatedAt = time.Now()
	expectedRentals[1].ID = 4
	expectedRentals[1].CreatedAt = time.Now()
	expectedRentals[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"rental_id", "rental_date", "inventory_id", "customer_id", "return_date", "staff_id",
	})
	for _, rental := range expectedRentals {
		rows.AddRow(
			rental.ID, rental.CreatedAt, rental.UpdatedAt, nil,
			rental.RentalID, rental.RentalDate, rental.InventoryID,
			rental.CustomerID, rental.ReturnDate, rental.StaffID,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rental` WHERE return_date IS NULL AND `rental`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	rentals, err := repo.FindNotReturned(context.Background())
	if err != nil {
		t.Errorf("Error finding not returned rentals: %v", err)
	}

	if len(rentals) != len(expectedRentals) {
		t.Errorf("Expected %d rentals, got %d", len(expectedRentals), len(rentals))
	}

	for _, rental := range rentals {
		if rental.ReturnDate != nil {
			t.Errorf("Expected ReturnDate to be nil, got %v", rental.ReturnDate)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
