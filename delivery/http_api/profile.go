package http_api

import (
	"sync"

	"github.com/Eco-Sort/eco_sort_backend/library/fiber_response"
	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/gofiber/fiber/v2"
)

type httpProfileApiDelivery struct {
}

func NewProfileHttpApiDelivery(app fiber.Router, middlewares ...func(*fiber.Ctx) error) {
	handler := httpProfileApiDelivery{}
	group := app.Group("profile", middlewares...)

	group.Get("/", handler.GetProfile)
}

func (h *httpProfileApiDelivery) GetProfile(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()
	userId := middleware.GetUserId(ctx)
	role := middleware.GetRole(ctx)
	return fiber_response.ReturnStatusOk(ctx, "Success", map[string]any{
		"user_id": userId,
		"role":    role,
	})
}
