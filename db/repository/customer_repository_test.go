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

func setupCustomerTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, CustomerRepository, func()) {
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

	repo := NewCustomerRepository(gormDB)

	return db, mock, repo, func() {
		db.Close()
	}
}

func TestCustomerRepository_Create(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	customer := &entity.Customer{
		StoreID:    1,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Address:    "123 Main St",
		Address2:   "Apt 4B",
		District:   "Downtown",
		City:       "Anytown",
		Country:    "USA",
		PostalCode: "12345",
		Phone:      "555-1234",
		Active:     true,
		CreateDate: time.Now(),
	}

	// Expect the INSERT query
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `customer`")).
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			customer.StoreID,
			customer.FirstName,
			customer.LastName,
			customer.Email,
			customer.Address,
			customer.Address2,
			customer.District,
			customer.City,
			customer.Country,
			customer.PostalCode,
			customer.Phone,
			customer.Active,
			customer.CreateDate,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Create(context.Background(), customer)
	if err != nil {
		t.Errorf("Error creating customer: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindByID(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customer
	expectedCustomer := entity.Customer{
		CustomerID: 1,
		StoreID:    1,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Address:    "123 Main St",
		Address2:   "Apt 4B",
		District:   "Downtown",
		City:       "Anytown",
		Country:    "USA",
		PostalCode: "12345",
		Phone:      "555-1234",
		Active:     true,
		CreateDate: time.Now(),
	}
	expectedCustomer.ID = 1
	expectedCustomer.CreatedAt = time.Now()
	expectedCustomer.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	}).
		AddRow(
			expectedCustomer.ID, expectedCustomer.CreatedAt, expectedCustomer.UpdatedAt, nil,
			expectedCustomer.CustomerID, expectedCustomer.StoreID, expectedCustomer.FirstName, expectedCustomer.LastName, expectedCustomer.Email,
			expectedCustomer.Address, expectedCustomer.Address2, expectedCustomer.District, expectedCustomer.City, expectedCustomer.Country,
			expectedCustomer.PostalCode, expectedCustomer.Phone, expectedCustomer.Active, expectedCustomer.CreateDate,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE `customer`.`id` = ? AND `customer`.`deleted_at` IS NULL ORDER BY `customer`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	customer, err := repo.FindByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding customer: %v", err)
	}

	if customer.CustomerID != expectedCustomer.CustomerID {
		t.Errorf("Expected CustomerID %d, got %d", expectedCustomer.CustomerID, customer.CustomerID)
	}

	if customer.FirstName != expectedCustomer.FirstName {
		t.Errorf("Expected FirstName %s, got %s", expectedCustomer.FirstName, customer.FirstName)
	}

	if customer.LastName != expectedCustomer.LastName {
		t.Errorf("Expected LastName %s, got %s", expectedCustomer.LastName, customer.LastName)
	}

	if customer.Email != expectedCustomer.Email {
		t.Errorf("Expected Email %s, got %s", expectedCustomer.Email, customer.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindAll(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customers
	expectedCustomers := []entity.Customer{
		{
			CustomerID: 1,
			StoreID:    1,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Address:    "123 Main St",
			Address2:   "Apt 4B",
			District:   "Downtown",
			City:       "Anytown",
			Country:    "USA",
			PostalCode: "12345",
			Phone:      "555-1234",
			Active:     true,
			CreateDate: time.Now(),
		},
		{
			CustomerID: 2,
			StoreID:    1,
			FirstName:  "Jane",
			LastName:   "Smith",
			Email:      "jane.smith@example.com",
			Address:    "456 Oak Ave",
			Address2:   "",
			District:   "Uptown",
			City:       "Othertown",
			Country:    "USA",
			PostalCode: "67890",
			Phone:      "555-5678",
			Active:     true,
			CreateDate: time.Now(),
		},
	}
	expectedCustomers[0].ID = 1
	expectedCustomers[0].CreatedAt = time.Now()
	expectedCustomers[0].UpdatedAt = time.Now()
	expectedCustomers[1].ID = 2
	expectedCustomers[1].CreatedAt = time.Now()
	expectedCustomers[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	})
	for _, customer := range expectedCustomers {
		rows.AddRow(
			customer.ID, customer.CreatedAt, customer.UpdatedAt, nil,
			customer.CustomerID, customer.StoreID, customer.FirstName, customer.LastName, customer.Email,
			customer.Address, customer.Address2, customer.District, customer.City, customer.Country,
			customer.PostalCode, customer.Phone, customer.Active, customer.CreateDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE `customer`.`deleted_at` IS NULL")).
		WillReturnRows(rows)

	customers, err := repo.FindAll(context.Background())
	if err != nil {
		t.Errorf("Error finding customers: %v", err)
	}

	if len(customers) != len(expectedCustomers) {
		t.Errorf("Expected %d customers, got %d", len(expectedCustomers), len(customers))
	}

	for i, customer := range customers {
		if customer.CustomerID != expectedCustomers[i].CustomerID {
			t.Errorf("Expected CustomerID %d, got %d", expectedCustomers[i].CustomerID, customer.CustomerID)
		}
		if customer.FirstName != expectedCustomers[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", expectedCustomers[i].FirstName, customer.FirstName)
		}
		if customer.LastName != expectedCustomers[i].LastName {
			t.Errorf("Expected LastName %s, got %s", expectedCustomers[i].LastName, customer.LastName)
		}
		if customer.Email != expectedCustomers[i].Email {
			t.Errorf("Expected Email %s, got %s", expectedCustomers[i].Email, customer.Email)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_Update(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	customer := &entity.Customer{
		CustomerID: 1,
		StoreID:    1,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Address:    "123 Main St",
		Address2:   "Apt 4B",
		District:   "Downtown",
		City:       "Anytown",
		Country:    "USA",
		PostalCode: "12345",
		Phone:      "555-1234",
		Active:     true,
		CreateDate: time.Now(),
	}
	customer.ID = 1

	// Expect the UPDATE query
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // CreatedAt
			sqlmock.AnyArg(), // UpdatedAt
			sqlmock.AnyArg(), // DeletedAt
			customer.StoreID,
			customer.FirstName,
			customer.LastName,
			customer.Email,
			customer.Address,
			customer.Address2,
			customer.District,
			customer.City,
			customer.Country,
			customer.PostalCode,
			customer.Phone,
			customer.Active,
			customer.CreateDate,
			customer.ID,
			customer.CustomerID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(context.Background(), customer)
	if err != nil {
		t.Errorf("Error updating customer: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_Delete(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	customer := &entity.Customer{
		CustomerID: 1,
		StoreID:    1,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
	}
	customer.ID = 1

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE").
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			customer.ID,
			customer.CustomerID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete(context.Background(), customer)
	if err != nil {
		t.Errorf("Error deleting customer: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_DeleteByID(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Expect the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `customer` SET")).
		WithArgs(
			sqlmock.AnyArg(), // DeletedAt
			1,                // ID
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.DeleteByID(context.Background(), 1)
	if err != nil {
		t.Errorf("Error deleting customer by ID: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindByName(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customers
	expectedCustomers := []entity.Customer{
		{
			CustomerID: 1,
			StoreID:    1,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Address:    "123 Main St",
			Address2:   "Apt 4B",
			District:   "Downtown",
			City:       "Anytown",
			Country:    "USA",
			PostalCode: "12345",
			Phone:      "555-1234",
			Active:     true,
			CreateDate: time.Now(),
		},
	}
	expectedCustomers[0].ID = 1
	expectedCustomers[0].CreatedAt = time.Now()
	expectedCustomers[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	})
	for _, customer := range expectedCustomers {
		rows.AddRow(
			customer.ID, customer.CreatedAt, customer.UpdatedAt, nil,
			customer.CustomerID, customer.StoreID, customer.FirstName, customer.LastName, customer.Email,
			customer.Address, customer.Address2, customer.District, customer.City, customer.Country,
			customer.PostalCode, customer.Phone, customer.Active, customer.CreateDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE (first_name LIKE ? OR last_name LIKE ?) AND `customer`.`deleted_at` IS NULL")).
		WithArgs("%John%", "%John%").
		WillReturnRows(rows)

	customers, err := repo.FindByName(context.Background(), "John")
	if err != nil {
		t.Errorf("Error finding customers by name: %v", err)
	}

	if len(customers) != len(expectedCustomers) {
		t.Errorf("Expected %d customers, got %d", len(expectedCustomers), len(customers))
	}

	for i, customer := range customers {
		if customer.FirstName != expectedCustomers[i].FirstName {
			t.Errorf("Expected FirstName %s, got %s", expectedCustomers[i].FirstName, customer.FirstName)
		}
		if customer.LastName != expectedCustomers[i].LastName {
			t.Errorf("Expected LastName %s, got %s", expectedCustomers[i].LastName, customer.LastName)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindByEmail(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customer
	expectedCustomer := entity.Customer{
		CustomerID: 1,
		StoreID:    1,
		FirstName:  "John",
		LastName:   "Doe",
		Email:      "john.doe@example.com",
		Address:    "123 Main St",
		Address2:   "Apt 4B",
		District:   "Downtown",
		City:       "Anytown",
		Country:    "USA",
		PostalCode: "12345",
		Phone:      "555-1234",
		Active:     true,
		CreateDate: time.Now(),
	}
	expectedCustomer.ID = 1
	expectedCustomer.CreatedAt = time.Now()
	expectedCustomer.UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	}).
		AddRow(
			expectedCustomer.ID, expectedCustomer.CreatedAt, expectedCustomer.UpdatedAt, nil,
			expectedCustomer.CustomerID, expectedCustomer.StoreID, expectedCustomer.FirstName, expectedCustomer.LastName, expectedCustomer.Email,
			expectedCustomer.Address, expectedCustomer.Address2, expectedCustomer.District, expectedCustomer.City, expectedCustomer.Country,
			expectedCustomer.PostalCode, expectedCustomer.Phone, expectedCustomer.Active, expectedCustomer.CreateDate,
		)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE email = ? AND `customer`.`deleted_at` IS NULL ORDER BY `customer`.`id` LIMIT ?")).
		WithArgs("john.doe@example.com", 1).
		WillReturnRows(rows)

	customer, err := repo.FindByEmail(context.Background(), "john.doe@example.com")
	if err != nil {
		t.Errorf("Error finding customer by email: %v", err)
	}

	if customer.Email != expectedCustomer.Email {
		t.Errorf("Expected Email %s, got %s", expectedCustomer.Email, customer.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindByStore(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customers
	expectedCustomers := []entity.Customer{
		{
			CustomerID: 1,
			StoreID:    1,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Address:    "123 Main St",
			Address2:   "Apt 4B",
			District:   "Downtown",
			City:       "Anytown",
			Country:    "USA",
			PostalCode: "12345",
			Phone:      "555-1234",
			Active:     true,
			CreateDate: time.Now(),
		},
		{
			CustomerID: 2,
			StoreID:    1,
			FirstName:  "Jane",
			LastName:   "Smith",
			Email:      "jane.smith@example.com",
			Address:    "456 Oak Ave",
			Address2:   "",
			District:   "Uptown",
			City:       "Othertown",
			Country:    "USA",
			PostalCode: "67890",
			Phone:      "555-5678",
			Active:     true,
			CreateDate: time.Now(),
		},
	}
	expectedCustomers[0].ID = 1
	expectedCustomers[0].CreatedAt = time.Now()
	expectedCustomers[0].UpdatedAt = time.Now()
	expectedCustomers[1].ID = 2
	expectedCustomers[1].CreatedAt = time.Now()
	expectedCustomers[1].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	})
	for _, customer := range expectedCustomers {
		rows.AddRow(
			customer.ID, customer.CreatedAt, customer.UpdatedAt, nil,
			customer.CustomerID, customer.StoreID, customer.FirstName, customer.LastName, customer.Email,
			customer.Address, customer.Address2, customer.District, customer.City, customer.Country,
			customer.PostalCode, customer.Phone, customer.Active, customer.CreateDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE store_id = ? AND `customer`.`deleted_at` IS NULL")).
		WithArgs(uint(1)).
		WillReturnRows(rows)

	customers, err := repo.FindByStore(context.Background(), 1)
	if err != nil {
		t.Errorf("Error finding customers by store: %v", err)
	}

	if len(customers) != len(expectedCustomers) {
		t.Errorf("Expected %d customers, got %d", len(expectedCustomers), len(customers))
	}

	for i, customer := range customers {
		if customer.StoreID != expectedCustomers[i].StoreID {
			t.Errorf("Expected StoreID %d, got %d", expectedCustomers[i].StoreID, customer.StoreID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindActive(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customers
	expectedCustomers := []entity.Customer{
		{
			CustomerID: 1,
			StoreID:    1,
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john.doe@example.com",
			Address:    "123 Main St",
			Address2:   "Apt 4B",
			District:   "Downtown",
			City:       "Anytown",
			Country:    "USA",
			PostalCode: "12345",
			Phone:      "555-1234",
			Active:     true,
			CreateDate: time.Now(),
		},
	}
	expectedCustomers[0].ID = 1
	expectedCustomers[0].CreatedAt = time.Now()
	expectedCustomers[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	})
	for _, customer := range expectedCustomers {
		rows.AddRow(
			customer.ID, customer.CreatedAt, customer.UpdatedAt, nil,
			customer.CustomerID, customer.StoreID, customer.FirstName, customer.LastName, customer.Email,
			customer.Address, customer.Address2, customer.District, customer.City, customer.Country,
			customer.PostalCode, customer.Phone, customer.Active, customer.CreateDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE active = ? AND `customer`.`deleted_at` IS NULL")).
		WithArgs(true).
		WillReturnRows(rows)

	customers, err := repo.FindActive(context.Background())
	if err != nil {
		t.Errorf("Error finding active customers: %v", err)
	}

	if len(customers) != len(expectedCustomers) {
		t.Errorf("Expected %d customers, got %d", len(expectedCustomers), len(customers))
	}

	for _, customer := range customers {
		if !customer.Active {
			t.Errorf("Expected Active to be true, got false")
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindInactive(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Define expected customers
	expectedCustomers := []entity.Customer{
		{
			CustomerID: 3,
			StoreID:    2,
			FirstName:  "Bob",
			LastName:   "Johnson",
			Email:      "bob.johnson@example.com",
			Address:    "789 Pine St",
			Address2:   "",
			District:   "Midtown",
			City:       "Sometown",
			Country:    "USA",
			PostalCode: "54321",
			Phone:      "555-9012",
			Active:     false,
			CreateDate: time.Now(),
		},
	}
	expectedCustomers[0].ID = 3
	expectedCustomers[0].CreatedAt = time.Now()
	expectedCustomers[0].UpdatedAt = time.Now()

	// Expect the SELECT query
	rows := sqlmock.NewRows([]string{
		"id", "created_at", "updated_at", "deleted_at",
		"customer_id", "store_id", "first_name", "last_name", "email",
		"address", "address2", "district", "city", "country",
		"postal_code", "phone", "active", "create_date",
	})
	for _, customer := range expectedCustomers {
		rows.AddRow(
			customer.ID, customer.CreatedAt, customer.UpdatedAt, nil,
			customer.CustomerID, customer.StoreID, customer.FirstName, customer.LastName, customer.Email,
			customer.Address, customer.Address2, customer.District, customer.City, customer.Country,
			customer.PostalCode, customer.Phone, customer.Active, customer.CreateDate,
		)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE active = ? AND `customer`.`deleted_at` IS NULL")).
		WithArgs(false).
		WillReturnRows(rows)

	customers, err := repo.FindInactive(context.Background())
	if err != nil {
		t.Errorf("Error finding inactive customers: %v", err)
	}

	if len(customers) != len(expectedCustomers) {
		t.Errorf("Expected %d customers, got %d", len(expectedCustomers), len(customers))
	}

	for _, customer := range customers {
		if customer.Active {
			t.Errorf("Expected Active to be false, got true")
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestCustomerRepository_FindByID_NotFound(t *testing.T) {
	_, mock, repo, cleanup := setupCustomerTest(t)
	defer cleanup()

	// Expect the SELECT query with no results
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `customer` WHERE `customer`.`id` = ? AND `customer`.`deleted_at` IS NULL ORDER BY `customer`.`id` LIMIT ?")).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindByID(context.Background(), 999)
	if err == nil {
		t.Error("Expected error when customer not found, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
