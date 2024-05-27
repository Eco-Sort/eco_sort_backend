package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/config"
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/golang-jwt/jwt"
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

func (c *authService) ValidateToken(token string) (domain.TokenPayload, error) {
	tokenValidation, err := jwt.Parse(token, func(tokenValidation *jwt.Token) (interface{}, error) {
		if _, ok := tokenValidation.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tokenValidation.Header["alg"])
		}
		return config.JWTSecretKey(), nil
	})
	if err != nil {
		return domain.TokenPayload{}, err
	}

	if !tokenValidation.Valid {
		return domain.TokenPayload{}, fmt.Errorf("invalid token")
	}

	claims, ok := tokenValidation.Claims.(jwt.MapClaims)
	if !ok {
		return domain.TokenPayload{}, fmt.Errorf("could not claim token")
	}

	expiration, ok := claims["exp"].(float64)
	if !ok {
		return domain.TokenPayload{}, fmt.Errorf("could not claim exp")
	}

	expirationTime := time.Unix(int64(expiration), 0)

	if time.Now().After(expirationTime) {
		return domain.TokenPayload{}, fmt.Errorf("token expired")
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return domain.TokenPayload{}, fmt.Errorf("could not claim user_id")
	}

	return domain.TokenPayload{
		UserId: uint(userId),
		Exp:    int64(expiration),
	}, nil
}
