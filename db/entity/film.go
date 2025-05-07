package entity

import (
	"gorm.io/gorm"
)

// Film represents a movie in the DVD rental system
type Film struct {
	gorm.Model
	FilmID      uint   `gorm:"primaryKey;column:film_id;autoIncrement"`
	Title       string `gorm:"column:title;not null"`
	ReleaseYear int16  `gorm:"column:release_year;not null"`
	Length      int16  `gorm:"column:length;not null"`
	CategoryID  uint   `gorm:"column:category_id;not null"`

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID"`
	Actors   []*Actor `gorm:"many2many:film_actors;foreignKey:FilmID;joinForeignKey:film_id;References:ActorID;joinReferences:actor_id"`
}

// TableName overrides the table name
func (Film) TableName() string {
	return "film"
}
