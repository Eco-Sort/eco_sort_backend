package http_api

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

type httpUserApiDelivery struct {
}

func NewAdminUserHttpApiDelivery(app fiber.Router) {
	handler := httpUserApiDelivery{}
	group := app.Group("user")

	group.Get("/", handler.GetUser)
}

func (h *httpUserApiDelivery) GetUser(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()
	statusCode := fiber.StatusOK
	return ctx.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": "Work",
		"data":    "",
	})
}
