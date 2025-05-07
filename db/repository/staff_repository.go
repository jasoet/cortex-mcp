package repository

import (
	"CortexMCP/db/entity"
	"context"
	"gorm.io/gorm"
)

// StaffRepository is an interface for staff operations
type StaffRepository interface {
	Repository[entity.Staff]

	// FindByName finds staff by first or last name
	FindByName(ctx context.Context, name string) ([]entity.Staff, error)

	// FindByEmail finds a staff member by email
	FindByEmail(ctx context.Context, email string) (*entity.Staff, error)

	// FindByUsername finds a staff member by username
	FindByUsername(ctx context.Context, username string) (*entity.Staff, error)

	// FindByStore finds staff by store ID
	FindByStore(ctx context.Context, storeID uint) ([]entity.Staff, error)

	// FindActive finds active staff
	FindActive(ctx context.Context) ([]entity.Staff, error)

	// FindInactive finds inactive staff
	FindInactive(ctx context.Context) ([]entity.Staff, error)
}

// StaffRepositoryImpl is an implementation of StaffRepository
type StaffRepositoryImpl struct {
	BaseRepository[entity.Staff]
}

// NewStaffRepository creates a new StaffRepository
func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &StaffRepositoryImpl{
		BaseRepository: BaseRepository[entity.Staff]{
			DB: db,
		},
	}
}

// FindByName finds staff by first or last name
func (r *StaffRepositoryImpl) FindByName(ctx context.Context, name string) ([]entity.Staff, error) {
	var staff []entity.Staff
	if err := r.DB.WithContext(ctx).Where("first_name LIKE ? OR last_name LIKE ?", "%"+name+"%", "%"+name+"%").Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// FindByEmail finds a staff member by email
func (r *StaffRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// FindByUsername finds a staff member by username
func (r *StaffRepositoryImpl) FindByUsername(ctx context.Context, username string) (*entity.Staff, error) {
	var staff entity.Staff
	if err := r.DB.WithContext(ctx).Where("username = ?", username).First(&staff).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

// FindByStore finds staff by store ID
func (r *StaffRepositoryImpl) FindByStore(ctx context.Context, storeID uint) ([]entity.Staff, error) {
	var staff []entity.Staff
	if err := r.DB.WithContext(ctx).Where("store_id = ?", storeID).Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// FindActive finds active staff
func (r *StaffRepositoryImpl) FindActive(ctx context.Context) ([]entity.Staff, error) {
	var staff []entity.Staff
	if err := r.DB.WithContext(ctx).Where("active = ?", true).Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

// FindInactive finds inactive staff
func (r *StaffRepositoryImpl) FindInactive(ctx context.Context) ([]entity.Staff, error) {
	var staff []entity.Staff
	if err := r.DB.WithContext(ctx).Where("active = ?", false).Find(&staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}
