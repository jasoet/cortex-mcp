package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	Store     StoreRepository
	Staff     StaffRepository
	Customer  CustomerRepository
	Category  CategoryRepository
	Film      FilmRepository
	Actor     ActorRepository
	Inventory InventoryRepository
	Rental    RentalRepository
	Payment   PaymentRepository
}

// NewRepositories creates and initializes all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Store:     NewStoreRepository(db),
		Staff:     NewStaffRepository(db),
		Customer:  NewCustomerRepository(db),
		Category:  NewCategoryRepository(db),
		Film:      NewFilmRepository(db),
		Actor:     NewActorRepository(db),
		Inventory: NewInventoryRepository(db),
		Rental:    NewRentalRepository(db),
		Payment:   NewPaymentRepository(db),
	}
}
