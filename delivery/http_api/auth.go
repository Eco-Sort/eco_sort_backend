package http_api

import (
	"net/mail"
	"sync"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/fiber_response"
	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type httpAuthApiDelivery struct {
	authService domain.AuthService
}

func NewAuthHttpApiDelivery(app fiber.Router, authService domain.AuthService) {
	handler := httpAuthApiDelivery{
		authService: authService,
	}

	group := app.Group("auth")

	group.Post("/login", handler.AuthLogin)
	group.Post("/register", handler.AuthRegister)
}

func (h *httpAuthApiDelivery) AuthLogin(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()

	req := new(domain.AuthLoginRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parse body", err)
	}

	res, er := govalidator.ValidateStruct(req)
	if !res {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parses body", er)
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Email not valid", err)
	}

	authRes, err := h.authService.AuthenticateAdmin(req)
	if err != nil {
		if err.Error() == "passsword mismatch" {
			return fiber_response.ReturnStatusUnauthorized(ctx)
		}
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, err.Error(), err)
	}
	tokenExpire := time.Now().Add(time.Hour * 24).Unix()

	token, err := middleware.CreateToken(authRes.ID, authRes.Role, tokenExpire)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, err.Error(), err)
	}
	statusCode := fiber.StatusOK
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error":   false,
		"message": "success",
		"loginResult": map[string]any{
			"userId": authRes.ID,
			"name":   authRes.Username,
			"token":  token,
		},
	})
}
func (h *httpAuthApiDelivery) AuthRegister(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()

	req := new(domain.AuthRegisterRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parse body", err)
	}

	res, er := govalidator.ValidateStruct(req)
	if !res {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parses body", er)
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Email not valid", err)
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return fiber_response.ReturnStatusServerError(ctx, "Failed to hash password", err)
	}

	req.Password = string(hpass)

	_, err = h.authService.Register(*req)
	if err != nil {
		return fiber_response.ReturnStatusServerError(ctx, "Failed to register user", err)
	}

	statusCode := fiber.StatusCreated
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error":    false,
		"messages": "User Created",
	})
}
