package garbage

import (
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
)

type garbageService struct {
	contextTimeout    time.Duration
	garbageRepository domain.GarbageRepository
}

func NewGarbageService(t time.Duration, garbageRepository domain.GarbageRepository) domain.GarbageService {
	return &garbageService{
		contextTimeout:    t,
		garbageRepository: garbageRepository,
	}
}

func (c *garbageService) Create(req domain.Garbage) (uint, error) {
	res, err := c.garbageRepository.Create(req)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (c *garbageService) GetById(id uint) (*domain.Garbage, error) {
	res, err := c.garbageRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
