package domain

import "gorm.io/gorm"

type Garbage struct {
	gorm.Model
	UserID     uint      `gorm:"not null" json:"user_id"`
	CategoryID uint      `gorm:"not null" json:"category_id"`
	Image      string    `json:"image"`
	User       *User     `gorm:"foreignKey:UserID"`
	Category   *Category `gorm:"foreignKey:CategoryID"`
}
