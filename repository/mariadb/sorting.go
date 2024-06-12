package mariadb

import (
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"gorm.io/gorm"
)

type mariadbSortingRepository struct {
	mariadb *gorm.DB
}

func NewMariadbSortingRepository(client *gorm.DB) domain.SortingRepository {
	return &mariadbSortingRepository{
		mariadb: client,
	}
}

func (r *mariadbSortingRepository) Create(sorting domain.SortingRequest) (uint, error) {
	newSorting := &domain.Sorting{
		SortingBin:   sorting.SortingBin,
		Instructions: sorting.Instructions,
	}
	result := r.mariadb.Create(newSorting)
	if result.Error != nil {
		return 0, result.Error
	}
	return newSorting.ID, nil
}
