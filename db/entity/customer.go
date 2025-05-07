package entity

import (
	"gorm.io/gorm"
	"time"
)

// Customer represents a customer in the DVD rental system
type Customer struct {
	gorm.Model
	CustomerID uint      `gorm:"primaryKey;column:customer_id;autoIncrement"`
	StoreID    uint      `gorm:"column:store_id;not null"`
	FirstName  string    `gorm:"column:first_name;not null"`
	LastName   string    `gorm:"column:last_name;not null"`
	Email      string    `gorm:"column:email;not null;unique"`
	Address    string    `gorm:"column:address;not null"`
	Address2   string    `gorm:"column:address2"`
	District   string    `gorm:"column:district;not null"`
	City       string    `gorm:"column:city;not null"`
	Country    string    `gorm:"column:country;not null"`
	PostalCode string    `gorm:"column:postal_code;not null"`
	Phone      string    `gorm:"column:phone;not null"`
	Active     bool      `gorm:"column:active;not null;default:true"`
	CreateDate time.Time `gorm:"column:create_date;not null;default:CURRENT_DATE"`

	// Relationships
	Store Store `gorm:"foreignKey:StoreID"`
}

// TableName overrides the table name
func (Customer) TableName() string {
	return "customer"
}
