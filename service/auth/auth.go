package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/config"
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	contextTimeout time.Duration
	userRepo       domain.UserRepository
}

func NewAuthService(t time.Duration, userRepo domain.UserRepository) domain.AuthService {
	return &authService{
		contextTimeout: t,
		userRepo:       userRepo,
	}
}

func (c *authService) AuthenticateAdmin(user *domain.AuthLoginRequest) (domain.User, error) {
	_, cancel := context.WithTimeout(context.Background(), c.contextTimeout)
	defer cancel()

	res, err := c.userRepo.GetByUsername(user.Username)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, errors.New("passsword mismatch")
	}
	return res, nil
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

	role, ok := claims["role"].(string)
	if !ok {
		return domain.TokenPayload{}, fmt.Errorf("could not claim role")
	}

	return domain.TokenPayload{
		UserId: uint(userId),
		Role:   domain.Role(role),
		Exp:    int64(expiration),
	}, nil
}

func (c *authService) Register(req domain.AuthRegisterRequest) (uint, error) {
	_, err := c.userRepo.GetByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	res, err := c.userRepo.Create(req)
	if err != nil {
		return 0, err
	}
	return res, nil
}
