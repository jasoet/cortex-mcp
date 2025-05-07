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

func setupPaymentTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, PaymentRepository, func()) {
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

	repo := NewPaymentRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestPaymentRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	paymentDate := time.Now()
	payment := &entity.Payment{
		CustomerID:  1,
		StaffID:     1,
		RentalID:    1,
		Amount:      9.99,
		PaymentDate: paymentDate,
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `payment`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			payment.PaymentID,
			payment.CustomerID,
			payment.StaffID,
			payment.RentalID,
			payment.Amount,
			payment.PaymentDate,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), payment)
	if err != nil {
		t.Errorf("Error creating payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected payment
	paymentDate := time.Now()
	expectedPayment := entity.Payment{
		PaymentID:   1,
		CustomerID:  1,
		StaffID:     1,
		RentalID:    1,
		Amount:      9.99,
		PaymentDate: paymentDate,
	}
	expectedPayment.ID = 1
	expectedPayment.CreatedAt = time.Now()
	expectedPayment.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	}).
		AddRow(
			expectedPayment.ID, expectedPayment.CreatedAt, expectedPayment.UpdatedAt, nil,
			expectedPayment.PaymentID, expectedPayment.CustomerID, expectedPayment.StaffID,
			expectedPayment.RentalID, expectedPayment.Amount, expectedPayment.PaymentDate,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE `payment`.`id` = ? AND `payment`.`deleted_at` IS NULL ORDER BY `payment`.`id` LIMIT 1")).
		WithArgs(1).
		WillReturnRows(rows)

	payment, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding payment: %v", err)
	}

	if payment.PaymentID != expectedPayment.PaymentID {
		t.Errorf("Expected PaymentID %d, got %d", expectedPayment.PaymentID, payment.PaymentID)
	}

	if payment.CustomerID != expectedPayment.CustomerID {
		t.Errorf("Expected CustomerID %d, got %d", expectedPayment.CustomerID, payment.CustomerID)
	}

	if payment.StaffID != expectedPayment.StaffID {
		t.Errorf("Expected StaffID %d, got %d", expectedPayment.StaffID, payment.StaffID)
	}

	if payment.RentalID != expectedPayment.RentalID {
		t.Errorf("Expected RentalID %d, got %d", expectedPayment.RentalID, payment.RentalID)
	}

	if payment.Amount != expectedPayment.Amount {
		t.Errorf("Expected Amount %f, got %f", expectedPayment.Amount, payment.Amount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected payments
	paymentDate1 := time.Now()
	paymentDate2 := time.Now().Add(time.Hour)
	expectedPayments := []entity.Payment{
		{
			PaymentID:   1,
			CustomerID:  1,
			StaffID:     1,
			RentalID:    1,
			Amount:      9.99,
			PaymentDate: paymentDate1,
		},
		{
			PaymentID:   2,
			CustomerID:  2,
			StaffID:     1,
			RentalID:    2,
			Amount:      5.99,
			PaymentDate: paymentDate2,
		},
	}
	expectedPayments[0].ID = 1
	expectedPayments[0].CreatedAt = time.Now()
	expectedPayments[0].UpdatedAt = time.Now()
	expectedPayments[1].ID = 2
	expectedPayments[1].CreatedAt = time.Now()
	expectedPayments[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	})
	for _, payment := range expectedPayments {
		rows.AddRow(
			payment.ID, payment.CreatedAt, payment.UpdatedAt, nil,
			payment.PaymentID, payment.CustomerID, payment.StaffID,
			payment.RentalID, payment.Amount, payment.PaymentDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE `payment`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	payments, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding payments: %v", err)
	}

	if len(payments) != len(expectedPayments) {
		t.Errorf("Expected %d payments, got %d", len(expectedPayments), len(payments))
	}

	for i, payment := range payments {
		if payment.PaymentID != expectedPayments[i].PaymentID {
			t.Errorf("Expected PaymentID %d, got %d", expectedPayments[i].PaymentID, payment.PaymentID)
		}
		if payment.CustomerID != expectedPayments[i].CustomerID {
			t.Errorf("Expected CustomerID %d, got %d", expectedPayments[i].CustomerID, payment.CustomerID)
		}
		if payment.StaffID != expectedPayments[i].StaffID {
			t.Errorf("Expected StaffID %d, got %d", expectedPayments[i].StaffID, payment.StaffID)
		}
		if payment.RentalID != expectedPayments[i].RentalID {
			t.Errorf("Expected RentalID %d, got %d", expectedPayments[i].RentalID, payment.RentalID)
		}
		if payment.Amount != expectedPayments[i].Amount {
			t.Errorf("Expected Amount %f, got %f", expectedPayments[i].Amount, payment.Amount)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	paymentDate := time.Now()
	payment := &entity.Payment{
		PaymentID:   1,
		CustomerID:  1,
		StaffID:     1,
		RentalID:    1,
		Amount:      9.99,
		PaymentDate: paymentDate,
	}
	payment.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			payment.PaymentID,
			payment.CustomerID,
			payment.StaffID,
			payment.RentalID,
			payment.Amount,
			payment.PaymentDate,
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), payment)
	if err != nil {
		t.Errorf("Error updating payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	payment := &entity.Payment{
		PaymentID:  1,
		CustomerID: 1,
		StaffID:    1,
		RentalID:   1,
	}
	payment.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `payment` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			payment.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), payment)
	if err != nil {
		t.Errorf("Error deleting payment: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `payment` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			1,                // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error deleting payment by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByCustomer(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected payments
	paymentDate1 := time.Now()
	paymentDate2 := time.Now().Add(time.Hour)
	expectedPayments := []entity.Payment{
		{
			PaymentID:   1,
			CustomerID:  1,
			StaffID:     1,
			RentalID:    1,
			Amount:      9.99,
			PaymentDate: paymentDate1,
		},
		{
			PaymentID:   2,
			CustomerID:  1,
			StaffID:     2,
			RentalID:    2,
			Amount:      5.99,
			PaymentDate: paymentDate2,
		},
	}
	expectedPayments[0].ID = 1
	expectedPayments[0].CreatedAt = time.Now()
	expectedPayments[0].UpdatedAt = time.Now()
	expectedPayments[1].ID = 2
	expectedPayments[1].CreatedAt = time.Now()
	expectedPayments[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	})
	for _, payment := range expectedPayments {
		rows.AddRow(
			payment.ID, payment.CreatedAt, payment.UpdatedAt, nil,
			payment.PaymentID, payment.CustomerID, payment.StaffID,
			payment.RentalID, payment.Amount, payment.PaymentDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE customer_id = ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	payments, err := repo.FindByCustomer(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding payments by customer: %v", err)
	}

	if len(payments) != len(expectedPayments) {
		t.Errorf("Expected %d payments, got %d", len(expectedPayments), len(payments))
	}

	for _, payment := range payments {
		if payment.CustomerID != 1 {
			t.Errorf("Expected CustomerID 1, got %d", payment.CustomerID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByStaff(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected payments
	paymentDate1 := time.Now()
	paymentDate2 := time.Now().Add(time.Hour)
	expectedPayments := []entity.Payment{
		{
			PaymentID:   1,
			CustomerID:  1,
			StaffID:     1,
			RentalID:    1,
			Amount:      9.99,
			PaymentDate: paymentDate1,
		},
		{
			PaymentID:   3,
			CustomerID:  2,
			StaffID:     1,
			RentalID:    3,
			Amount:      7.99,
			PaymentDate: paymentDate2,
		},
	}
	expectedPayments[0].ID = 1
	expectedPayments[0].CreatedAt = time.Now()
	expectedPayments[0].UpdatedAt = time.Now()
	expectedPayments[1].ID = 3
	expectedPayments[1].CreatedAt = time.Now()
	expectedPayments[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	})
	for _, payment := range expectedPayments {
		rows.AddRow(
			payment.ID, payment.CreatedAt, payment.UpdatedAt, nil,
			payment.PaymentID, payment.CustomerID, payment.StaffID,
			payment.RentalID, payment.Amount, payment.PaymentDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE staff_id = ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	payments, err := repo.FindByStaff(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding payments by staff: %v", err)
	}

	if len(payments) != len(expectedPayments) {
		t.Errorf("Expected %d payments, got %d", len(expectedPayments), len(payments))
	}

	for _, payment := range payments {
		if payment.StaffID != 1 {
			t.Errorf("Expected StaffID 1, got %d", payment.StaffID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByRental(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected payment
	paymentDate := time.Now()
	expectedPayment := entity.Payment{
		PaymentID:   1,
		CustomerID:  1,
		StaffID:     1,
		RentalID:    1,
		Amount:      9.99,
		PaymentDate: paymentDate,
	}
	expectedPayment.ID = 1
	expectedPayment.CreatedAt = time.Now()
	expectedPayment.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	}).
		AddRow(
			expectedPayment.ID, expectedPayment.CreatedAt, expectedPayment.UpdatedAt, nil,
			expectedPayment.PaymentID, expectedPayment.CustomerID, expectedPayment.StaffID,
			expectedPayment.RentalID, expectedPayment.Amount, expectedPayment.PaymentDate,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE rental_id = ? AND `payment`.`deleted_at` IS NULL ORDER BY `payment`.`id` LIMIT 1")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	payment, err := repo.FindByRental(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding payment by rental: %v", err)
	}

	if payment.RentalID != 1 {
		t.Errorf("Expected RentalID 1, got %d", payment.RentalID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByDateRange(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define date range
	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 1, 31, 23, 59, 59, 0, time.UTC)

	// Define expected payments
	paymentDate1 := time.Date(2023, 1, 15, 12, 0, 0, 0, time.UTC)
	paymentDate2 := time.Date(2023, 1, 20, 14, 30, 0, 0, time.UTC)
	expectedPayments := []entity.Payment{
		{
			PaymentID:   1,
			CustomerID:  1,
			StaffID:     1,
			RentalID:    1,
			Amount:      9.99,
			PaymentDate: paymentDate1,
		},
		{
			PaymentID:   2,
			CustomerID:  2,
			StaffID:     1,
			RentalID:    2,
			Amount:      5.99,
			PaymentDate: paymentDate2,
		},
	}
	expectedPayments[0].ID = 1
	expectedPayments[0].CreatedAt = time.Now()
	expectedPayments[0].UpdatedAt = time.Now()
	expectedPayments[1].ID = 2
	expectedPayments[1].CreatedAt = time.Now()
	expectedPayments[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	})
	for _, payment := range expectedPayments {
		rows.AddRow(
			payment.ID, payment.CreatedAt, payment.UpdatedAt, nil,
			payment.PaymentID, payment.CustomerID, payment.StaffID,
			payment.RentalID, payment.Amount, payment.PaymentDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE payment_date BETWEEN ? AND ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(startDate, endDate).
		WillReturnRows(rows)

	payments, err := repo.FindByDateRange(context.Background(), startDate, endDate)
	if err != nil {
		t.Errorf("Error finding payments by date range: %v", err)
	}

	if len(payments) != len(expectedPayments) {
		t.Errorf("Expected %d payments, got %d", len(expectedPayments), len(payments))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_FindByAmountRange(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define amount range
	minAmount := 5.0
	maxAmount := 10.0

	// Define expected payments
	paymentDate1 := time.Now()
	paymentDate2 := time.Now().Add(time.Hour)
	expectedPayments := []entity.Payment{
		{
			PaymentID:   1,
			CustomerID:  1,
			StaffID:     1,
			RentalID:    1,
			Amount:      9.99,
			PaymentDate: paymentDate1,
		},
		{
			PaymentID:   2,
			CustomerID:  2,
			StaffID:     1,
			RentalID:    2,
			Amount:      5.99,
			PaymentDate: paymentDate2,
		},
	}
	expectedPayments[0].ID = 1
	expectedPayments[0].CreatedAt = time.Now()
	expectedPayments[0].UpdatedAt = time.Now()
	expectedPayments[1].ID = 2
	expectedPayments[1].CreatedAt = time.Now()
	expectedPayments[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"payment_id", "customer_id", "staff_id", "rental_id", "amount", "payment_date",
	})
	for _, payment := range expectedPayments {
		rows.AddRow(
			payment.ID, payment.CreatedAt, payment.UpdatedAt, nil,
			payment.PaymentID, payment.CustomerID, payment.StaffID,
			payment.RentalID, payment.Amount, payment.PaymentDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `payment` WHERE amount BETWEEN ? AND ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(minAmount, maxAmount).
		WillReturnRows(rows)

	payments, err := repo.FindByAmountRange(context.Background(), minAmount, maxAmount)
	if err != nil {
		t.Errorf("Error finding payments by amount range: %v", err)
	}

	if len(payments) != len(expectedPayments) {
		t.Errorf("Expected %d payments, got %d", len(expectedPayments), len(payments))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_GetTotalPaymentsByCustomer(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected total
	expectedTotal := 15.98

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"SUM(amount)"}).
		AddRow(expectedTotal)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(amount) FROM `payment` WHERE customer_id = ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	total, err := repo.GetTotalPaymentsByCustomer(context.Background(), 1)
	if err != nil {
		t.Errorf("Error getting total payments by customer: %v", err)
	}

	if total != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestPaymentRepository_GetTotalPaymentsByStore(t *testing.T) {
	_, mock, repo, cleanup := setupPaymentTest(t)
	defer cleanup()

	// Define expected total
	expectedTotal := 25.97

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{"SUM(payment.amount)"}).
		AddRow(expectedTotal)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(payment.amount) FROM `payment` JOIN staff ON payment.staff_id = staff.staff_id WHERE staff.store_id = ? AND `payment`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	total, err := repo.GetTotalPaymentsByStore(context.Background(), 1)
	if err != nil {
		t.Errorf("Error getting total payments by store: %v", err)
	}

	if total != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, total)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
