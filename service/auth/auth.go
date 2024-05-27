package auth

import (
	"context"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
)

type authService struct {
	contextTimeout time.Duration
}

func NewAuthService(t time.Duration) domain.AuthService {
	return &authService{
		contextTimeout: t,
	}
}

func (c *authService) AuthenticateAdmin(user *domain.AuthLoginRequest) (bool, error) {
	_, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	if user.Username == "admin" && user.Password == "admin" {
		return true, nil
	}

	return false, nil
}
