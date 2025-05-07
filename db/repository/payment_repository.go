package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
	"time"
)

// PaymentRepository is an interface for payment operations
type PaymentRepository interface {
	Repository[entity.Payment]

	// FindByCustomer finds payments by customer ID
	FindByCustomer(ctx context.Context, customerID uint) ([]entity.Payment, error)

	// FindByStaff finds payments by staff ID
	FindByStaff(ctx context.Context, staffID uint) ([]entity.Payment, error)

	// FindByRental finds a payment by rental ID
	FindByRental(ctx context.Context, rentalID uint) (*entity.Payment, error)

	// FindByDateRange finds payments within a date range
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Payment, error)

	// FindByAmountRange finds payments within an amount range
	FindByAmountRange(ctx context.Context, minAmount, maxAmount float64) ([]entity.Payment, error)

	// GetTotalPaymentsByCustomer gets the total amount of payments by customer ID
	GetTotalPaymentsByCustomer(ctx context.Context, customerID uint) (float64, error)

	// GetTotalPaymentsByStore gets the total amount of payments by store ID
	GetTotalPaymentsByStore(ctx context.Context, storeID uint) (float64, error)
}

// PaymentRepositoryImpl is an implementation of PaymentRepository
type PaymentRepositoryImpl struct {
	BaseRepository[entity.Payment]
}

// NewPaymentRepository creates a new PaymentRepository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{
		BaseRepository: BaseRepository[entity.Payment]{
			DB: db,
		},
	}
}

// FindByCustomer finds payments by customer ID
func (r *PaymentRepositoryImpl) FindByCustomer(ctx context.Context, customerID uint) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := r.DB.WithContext(ctx).Where("customer_id = ?", customerID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindByStaff finds payments by staff ID
func (r *PaymentRepositoryImpl) FindByStaff(ctx context.Context, staffID uint) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := r.DB.WithContext(ctx).Where("staff_id = ?", staffID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindByRental finds a payment by rental ID
func (r *PaymentRepositoryImpl) FindByRental(ctx context.Context, rentalID uint) (*entity.Payment, error) {
	var payment entity.Payment
	if err := r.DB.WithContext(ctx).Where("rental_id = ?", rentalID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByDateRange finds payments within a date range
func (r *PaymentRepositoryImpl) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := r.DB.WithContext(ctx).Where("payment_date BETWEEN ? AND ?", startDate, endDate).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// FindByAmountRange finds payments within an amount range
func (r *PaymentRepositoryImpl) FindByAmountRange(ctx context.Context, minAmount, maxAmount float64) ([]entity.Payment, error) {
	var payments []entity.Payment
	if err := r.DB.WithContext(ctx).Where("amount BETWEEN ? AND ?", minAmount, maxAmount).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

// GetTotalPaymentsByCustomer gets the total amount of payments by customer ID
func (r *PaymentRepositoryImpl) GetTotalPaymentsByCustomer(ctx context.Context, customerID uint) (float64, error) {
	var total float64
	if err := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Select("SUM(amount)").
		Where("customer_id = ?", customerID).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// GetTotalPaymentsByStore gets the total amount of payments by store ID
func (r *PaymentRepositoryImpl) GetTotalPaymentsByStore(ctx context.Context, storeID uint) (float64, error) {
	var total float64
	if err := r.DB.WithContext(ctx).Model(&entity.Payment{}).
		Select("SUM(payment.amount)").
		Joins("JOIN staff ON payment.staff_id = staff.staff_id").
		Where("staff.store_id = ?", storeID).
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
