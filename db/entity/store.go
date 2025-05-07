package entity

import (
	"gorm.io/gorm"
)

// Store represents a store in the DVD rental system
type Store struct {
	gorm.Model
	StoreID    uint   `gorm:"primaryKey;column:store_id;autoIncrement"`
	StoreName  string `gorm:"column:store_name;not null"`
	Address    string `gorm:"column:address;not null"`
	Address2   string `gorm:"column:address2"`
	District   string `gorm:"column:district;not null"`
	City       string `gorm:"column:city;not null"`
	Country    string `gorm:"column:country;not null"`
	PostalCode string `gorm:"column:postal_code;not null"`
	Phone      string `gorm:"column:phone;not null"`

	// Relationships
}

// TableName overrides the table name
func (Store) TableName() string {
	return "store"
}
