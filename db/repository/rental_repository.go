package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
	"time"
)

// RentalRepository is an interface for rental operations
type RentalRepository interface {
	Repository[entity.Rental]

	// FindByCustomer finds rentals by customer ID
	FindByCustomer(ctx context.Context, customerID uint) ([]entity.Rental, error)

	// FindByStaff finds rentals by staff ID
	FindByStaff(ctx context.Context, staffID uint) ([]entity.Rental, error)

	// FindByInventory finds rentals by inventory ID
	FindByInventory(ctx context.Context, inventoryID uint) ([]entity.Rental, error)

	// FindByDateRange finds rentals within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Rental, error)

	// FindOverdue finds overdue rentals (no return date and rental date is older than specified days)
	FindOverdue(ctx context.Context, daysOverdue int) ([]entity.Rental, error)

	// FindReturned finds rentals that have been returned
	FindReturned(ctx context.Context) ([]entity.Rental, error)

	// FindNotReturned finds rentals that have not been returned
	FindNotReturned(ctx context.Context) ([]entity.Rental, error)
}

// RentalRepositoryImpl is an implementation of RentalRepository
type RentalRepositoryImpl struct {
	BaseRepository[entity.Rental]
}

// NewRentalRepository creates a new RentalRepository
func NewRentalRepository(db *gorm.DB) RentalRepository {
	return &RentalRepositoryImpl{
		BaseRepository: BaseRepository[entity.Rental]{
			DB: db,
		},
	}
}

// FindByCustomer finds rentals by customer ID
func (r *RentalRepositoryImpl) FindByCustomer(ctx context.Context, customerID uint) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("customer_id = ?", customerID).Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindByStaff finds rentals by staff ID
func (r *RentalRepositoryImpl) FindByStaff(ctx context.Context, staffID uint) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("staff_id = ?", staffID).Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindByInventory finds rentals by inventory ID
func (r *RentalRepositoryImpl) FindByInventory(ctx context.Context, inventoryID uint) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("inventory_id = ?", inventoryID).Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindByDateRange finds rentals within a date range
func (r *RentalRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("rental_date BETWEEN ? AND ?", startDate, endDate).Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindOverdue finds overdue rentals (no return date and rental date is older than specified days)
func (r *RentalRepositoryImpl) FindOverdue(ctx context.Context, daysOverdue int) ([]entity.Rental, error) {
	var rentals []entity.Rental
	overdueCutoff := time.Now().AddDate(0, 0, -daysOverdue)
	if err := r.DB.WithContext(ctx).
		Where("return_date IS NULL AND rental_date < ?", overdueCutoff).
		Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindReturned finds rentals that have been returned
func (r *RentalRepositoryImpl) FindReturned(ctx context.Context) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("return_date IS NOT NULL").Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}

// FindNotReturned finds rentals that have not been returned
func (r *RentalRepositoryImpl) FindNotReturned(ctx context.Context) ([]entity.Rental, error) {
	var rentals []entity.Rental
	if err := r.DB.WithContext(ctx).Where("return_date IS NULL").Find(&rentals).Error; err != nil {
		return nil, err
	}
	return rentals, nil
}
