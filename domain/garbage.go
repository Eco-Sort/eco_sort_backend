package domain

import "gorm.io/gorm"

type Garbage struct {
	gorm.Model
	UserID      uint          `gorm:"not null" json:"user_id"`
	Image       string        `json:"image"`
	ImageObject []ImageObject `json:"image_object"`
	User        *User         `gorm:"foreignKey:UserID"`
}

type ImageObject struct {
	gorm.Model
	CategoryID uint      `gorm:"not null" json:"category_id"`
	GarbageID  uint      `gorm:"not null" json:"garbage_id"`
	Confidence float64   `json:"confidence"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
}

type GarbageRepository interface {
	Create(category Garbage) (uint, error)
	GetById(id uint) (*Garbage, error)
}

type GarbageService interface {
	Create(req Garbage) (uint, error)
	GetById(id uint) (*Garbage, error)
}
