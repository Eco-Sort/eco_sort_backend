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

func (r *mariadbGarbageRepository) Create(garbage domain.Garbage) (uint, error) {
	newGarbage := garbage
	result := r.mariadb.Create(&newGarbage)
	if result.Error != nil {
		return 0, result.Error
	}
	for _, s := range newGarbage.ImageObject {
		s.GarbageID = newGarbage.ID
		res := r.mariadb.Omit("id").Create(&s)
		if res.Error != nil {
			return 0, res.Error
		}
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
