package http_api

import (
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
	return fiber_response.ReturnStatusOk(ctx, "Welcome", map[string]any{
		"token":   token,
		"user_id": authRes.ID,
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

	if req.Password != req.ReEnterPassword {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Password must be the same", nil)
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return fiber_response.ReturnStatusServerError(ctx, "Failed to hash password", err)
	}

	req.Password = string(hpass)

	result, err := h.authService.Register(*req)
	if err != nil {
		return fiber_response.ReturnStatusServerError(ctx, "Failed to register user", err)
	}

	return fiber_response.ReturnStatusCreated(ctx, "Success", map[string]any{
		"user_id": result,
	})
}
