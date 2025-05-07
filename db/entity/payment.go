package entity

import (
	"gorm.io/gorm"
	"time"
)

// Payment represents a payment for a rental in the DVD rental system
type Payment struct {
	gorm.Model
	PaymentID   uint      `gorm:"primaryKey;column:payment_id;autoIncrement"`
	CustomerID  uint      `gorm:"column:customer_id;not null"`
	StaffID     uint      `gorm:"column:staff_id;not null"`
	RentalID    uint      `gorm:"column:rental_id;not null;unique"`
	Amount      float64   `gorm:"column:amount;not null;type:numeric(5,2)"`
	PaymentDate time.Time `gorm:"column:payment_date;not null"`

	// Relationships
	Customer Customer `gorm:"foreignKey:CustomerID"`
	Staff    Staff    `gorm:"foreignKey:StaffID"`
	Rental   Rental   `gorm:"foreignKey:RentalID"`
}

// TableName overrides the table name
func (Payment) TableName() string {
	return "payment"
}
