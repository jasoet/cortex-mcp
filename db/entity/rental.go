package entity

import (
	"gorm.io/gorm"
	"time"
)

// Rental represents a film rental transaction in the DVD rental system
type Rental struct {
	gorm.Model
	RentalID    uint       `gorm:"primaryKey;column:rental_id;autoIncrement"`
	RentalDate  time.Time  `gorm:"column:rental_date;not null"`
	InventoryID uint       `gorm:"column:inventory_id;not null"`
	CustomerID  uint       `gorm:"column:customer_id;not null"`
	ReturnDate  *time.Time `gorm:"column:return_date"`
	StaffID     uint       `gorm:"column:staff_id;not null"`

	// Relationships
	Inventory Inventory `gorm:"foreignKey:InventoryID"`
	Customer  Customer  `gorm:"foreignKey:CustomerID"`
	Staff     Staff     `gorm:"foreignKey:StaffID"`
	Payment   *Payment  `gorm:"foreignKey:RentalID"`
}

// TableName overrides the table name
func (Rental) TableName() string {
	return "rental"
}
