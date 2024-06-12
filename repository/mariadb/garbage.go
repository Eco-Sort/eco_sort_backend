package mariadb

import (
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"gorm.io/gorm"
)

type mariadbGarbageRepository struct {
	mariadb *gorm.DB
}

func NewMariadbGarbageRepository(client *gorm.DB) domain.GarbageRepository {
	return &mariadbGarbageRepository{
		mariadb: client,
	}
}

func (r *mariadbGarbageRepository) Create(category domain.Garbage) (uint, error) {
	newGarbage := category
	for _, imageObject := range newGarbage.ImageObject {
		res := r.mariadb.Create(&imageObject)
		if res.Error != nil {
			return 0, res.Error
		}
	}
	result := r.mariadb.Create(&newGarbage)
	if result.Error != nil {
		return 0, result.Error
	}
	return newGarbage.ID, nil
}

func (r *mariadbGarbageRepository) GetById(id uint) (*domain.Garbage, error) {
	garbage := &domain.Garbage{}
	result := r.mariadb.Preload("ImageObject").Preload("ImageObject.Category").Preload("ImageObject.Category.Sorting").First(garbage, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return garbage, nil
}
