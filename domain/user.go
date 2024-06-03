package domain

import "gorm.io/gorm"

type Role string

const (
	Admin  Role = "admin"
	Client Role = "client"
	Guest  Role = "guest"
)

type User struct {
	gorm.Model
	Role     Role   `gorm:"not null" json:"role"`
	Username string `gorm:"not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}

type UserRepository interface {
	Create(user AuthRegisterRequest) (uint, error)
	GetByUserId(userId uint) (User, error)
	GetByUsername(username string) (User, error)
}
