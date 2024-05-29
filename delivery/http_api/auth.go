package http_api

import (
	"fmt"
	"sync"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/fiber_response"
	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
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
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, err.Error(), err)
	}
	if !authRes {
		return fiber_response.ReturnStatusUnauthorized(ctx)
	}
	tokenExpire := time.Now().Add(time.Hour * 24).Unix()
	//TODO:Changing userId
	token, err := middleware.CreateToken(1, domain.Admin, tokenExpire)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, err.Error(), err)
	}
	return fiber_response.ReturnStatusOk(ctx, fmt.Sprintf("Welcome %s", req.Username), map[string]any{
		"token":  token,
		"userId": 1,
	})
}
