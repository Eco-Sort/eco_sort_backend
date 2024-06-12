package mariadb

import (
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"gorm.io/gorm"
)

type mariadbCategoryRepository struct {
	mariadb *gorm.DB
}

func NewMariadbCategoryRepository(client *gorm.DB) domain.CategoryRepository {
	return &mariadbCategoryRepository{
		mariadb: client,
	}
}

func (r *mariadbCategoryRepository) Create(category domain.CategoryRequest) (uint, error) {
	newCategory := &domain.Category{
		Label:     category.Label,
		SortingID: category.SortingID,
	}
	result := r.mariadb.Create(newCategory)
	if result.Error != nil {
		return 0, result.Error
	}
	return newCategory.ID, nil
}
