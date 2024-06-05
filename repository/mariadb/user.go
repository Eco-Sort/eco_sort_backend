package mariadb

import (
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"gorm.io/gorm"
)

// var _userRepoName = "users"

type mariadbUserRepository struct {
	mariadb *gorm.DB
}

func NewMariadbUserRepository(client *gorm.DB) domain.UserRepository {
	return &mariadbUserRepository{
		mariadb: client,
	}
}

func (u *mariadbUserRepository) Create(user domain.AuthRegisterRequest) (uint, error) {
	newUser := &domain.User{
		Role:     domain.Client,
		Username: user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	result := u.mariadb.Create(newUser)
	if result.Error != nil {
		return 0, result.Error
	}
	return newUser.ID, nil
}
func (u *mariadbUserRepository) GetByUserId(userId uint) (domain.User, error) {
	var user domain.User
	result := u.mariadb.First(&user, userId)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (u *mariadbUserRepository) GetByEmail(email string) (domain.User, error) {
	var user domain.User
	result := u.mariadb.First(&user, "email = ?", email)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}
