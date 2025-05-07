package entity

import (
	"gorm.io/gorm"
)

// Category represents a film genre in the DVD rental system
type Category struct {
	gorm.Model
	CategoryID uint   `gorm:"primaryKey;column:category_id;autoIncrement"`
	Name       string `gorm:"column:name;not null;unique"`
}

// TableName overrides the table name
func (Category) TableName() string {
	return "category"
}
