package entity

import (
	"gorm.io/gorm"
	"time"
)

// Staff represents an employee in the DVD rental system
type Staff struct {
	gorm.Model
	StaffID    uint      `gorm:"primaryKey;column:staff_id;autoIncrement"`
	StoreID    uint      `gorm:"column:store_id;not null"`
	FirstName  string    `gorm:"column:first_name;not null"`
	LastName   string    `gorm:"column:last_name;not null"`
	Email      string    `gorm:"column:email;not null;unique"`
	Username   string    `gorm:"column:username;not null;unique"`
	Address    string    `gorm:"column:address;not null"`
	Address2   string    `gorm:"column:address2"`
	District   string    `gorm:"column:district;not null"`
	City       string    `gorm:"column:city;not null"`
	Country    string    `gorm:"column:country;not null"`
	PostalCode string    `gorm:"column:postal_code;not null"`
	Phone      string    `gorm:"column:phone;not null"`
	Active     bool      `gorm:"column:active;not null;default:true"`
	LastUpdate time.Time `gorm:"column:last_update;not null;default:CURRENT_TIMESTAMP"`

	// Relationships
	Store Store `gorm:"foreignKey:StoreID"`
}

// TableName overrides the table name
func (Staff) TableName() string {
	return "staff"
}
