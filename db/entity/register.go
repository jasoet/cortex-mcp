package entity

import (
	"gorm.io/gorm"
)

// RegisterEntities registers all entities with GORM
func RegisterEntities(db *gorm.DB) error {
	// Register all entities for auto-migration if needed
	if err := db.AutoMigrate(
		&Store{},
		&Staff{},
		&Customer{},
		&Category{},
		&Film{},
		&Actor{},
		&Inventory{},
		&Rental{},
		&Payment{},
	); err != nil {
		return err
	}

	return nil
}
