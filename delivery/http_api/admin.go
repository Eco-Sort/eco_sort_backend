package http_api

import (
	"sync"

	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/gofiber/fiber/v2"
)

type httpUserApiDelivery struct {
}

func NewAdminUserHttpApiDelivery(app fiber.Router, middlewares ...func(*fiber.Ctx) error) {
	handler := httpUserApiDelivery{}
	group := app.Group("user", middlewares...)

	group.Get("/", handler.GetUser)
}

func (h *httpUserApiDelivery) GetUser(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()
	userId := middleware.GetUserId(ctx)
	statusCode := fiber.StatusOK
	return ctx.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": "Work",
		"data":    userId,
	})
}
