package entity

import (
	"gorm.io/gorm"
)

// Actor represents an actor in the DVD rental system
type Actor struct {
	gorm.Model
	ActorID   uint   `gorm:"primaryKey;column:actor_id;autoIncrement"`
	FirstName string `gorm:"column:first_name;not null"`
	LastName  string `gorm:"column:last_name;not null"`

	// Relationships
}

// TableName overrides the table name
func (Actor) TableName() string {
	return "actor"
}
