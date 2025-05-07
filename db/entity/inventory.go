package entity

import (
	"gorm.io/gorm"
)

// Inventory represents a copy of a film in a store in the DVD rental system
type Inventory struct {
	gorm.Model
	InventoryID uint `gorm:"primaryKey;column:inventory_id;autoIncrement"`
	FilmID      uint `gorm:"column:film_id;not null"`
	StoreID     uint `gorm:"column:store_id;not null"`

	// Relationships
	Film  Film  `gorm:"foreignKey:FilmID"`
	Store Store `gorm:"foreignKey:StoreID"`
}

// TableName overrides the table name
func (Inventory) TableName() string {
	return "inventory"
}
