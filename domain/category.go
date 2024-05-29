package domain

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	SortingID uint     `gorm:"not null" json:"sorting_id"`
	Label     string   `gorm:"not null" json:"label"`
	Sorting   *Sorting `gorm:"foreignKey:SortingID"`
}
