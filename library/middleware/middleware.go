package middleware

import (
	"strings"

	"github.com/Eco-Sort/eco_sort_backend/config"
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/fiber_response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetUserId(ctx *fiber.Ctx) uint {
	return ctx.Locals("user_id").(uint)
}

func GetRole(ctx *fiber.Ctx) domain.Role {
	return ctx.Locals("role").(domain.Role)
}
func CreateToken(userId uint, role domain.Role, tokenExpire int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": int64(userId),
		"role":    role,
		"exp":     tokenExpire,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(config.JWTSecretKey())
	if err != nil {
		return "", err
	}

	return t, nil
}

func ValidateJWT(ctx *fiber.Ctx) error {
	authService := *ctx.Locals("auth_service").(*domain.AuthService)
	authHeaders := string(ctx.Request().Header.Peek("Authorization"))
	if !strings.Contains(authHeaders, "Bearer") {
		return fiber_response.ReturnStatusUnauthorized(ctx)
	}

	token := strings.Replace(authHeaders, "Bearer ", "", -1)
	if token == "Bearer" {
		return fiber_response.ReturnStatusUnauthorized(ctx)
	}

	res, err := authService.ValidateToken(token)
	if err != nil {
		return fiber_response.ReturnStatusUnauthorized(ctx)
	} else {
		ctx.Locals("role", res.Role)
		ctx.Locals("user_id", res.UserId)
		return ctx.Next()
	}
}
