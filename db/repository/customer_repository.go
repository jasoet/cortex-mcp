package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// CustomerRepository is an interface for customer operations
type CustomerRepository interface {
	Repository[entity.Customer]

	// FindByName finds customers by first or last name
	FindByName(ctx context.Context, name string) ([]entity.Customer, error)

	// FindByEmail finds a customer by email
	FindByEmail(ctx context.Context, email string) (*entity.Customer, error)

	// FindByStore finds customers by store ID
	FindByStore(ctx context.Context, storeID uint) ([]entity.Customer, error)

	// FindActive finds active customers
	FindActive(ctx context.Context) ([]entity.Customer, error)

	// FindInactive finds inactive customers
	FindInactive(ctx context.Context) ([]entity.Customer, error)
}

// CustomerRepositoryImpl is an implementation of CustomerRepository
type CustomerRepositoryImpl struct {
	BaseRepository[entity.Customer]
}

// NewCustomerRepository creates a new CustomerRepository
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{
		BaseRepository: BaseRepository[entity.Customer]{
			DB: db,
		},
	}
}

// FindByName finds customers by first or last name
func (r *CustomerRepositoryImpl) FindByName(ctx context.Context, name string) ([]entity.Customer, error) {
	var customers []entity.Customer
	if err := r.DB.WithContext(ctx).Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%").Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// FindByEmail finds a customer by email
func (r *CustomerRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	var customer entity.Customer
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// FindByStore finds customers by store ID
func (r *CustomerRepositoryImpl) FindByStore(ctx context.Context, storeID uint) ([]entity.Customer, error) {
	var customers []entity.Customer
	if err := r.DB.WithContext(ctx).Where("store_id = ?", storeID).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// FindActive finds active customers
func (r *CustomerRepositoryImpl) FindActive(ctx context.Context) ([]entity.Customer, error) {
	var customers []entity.Customer
	if err := r.DB.WithContext(ctx).Where("active = ?", true).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// FindInactive finds inactive customers
func (r *CustomerRepositoryImpl) FindInactive(ctx context.Context) ([]entity.Customer, error) {
	var customers []entity.Customer
	if err := r.DB.WithContext(ctx).Where("active = ?", false).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
