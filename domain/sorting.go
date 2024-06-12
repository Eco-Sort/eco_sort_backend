package domain

import "gorm.io/gorm"

type Sorting struct {
	gorm.Model
	SortingBin   string `gorm:"not null" json:"sorting_bin"`
	Instructions string `gorm:"not null" json:"instructions"`
}

type SortingRequest struct {
	SortingBin   string `json:"sorting_bin"`
	Instructions string `json:"instructions"`
}

type SortingRepository interface {
	Create(sorting SortingRequest) (uint, error)
}
