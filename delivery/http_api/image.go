package http_api

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/fiber_response"
	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
)

type httpImageApiDelivery struct {
	imageService domain.ImageService
}

func NewImageHttpApiDelivery(app fiber.Router, imageService domain.ImageService, middlewares ...func(*fiber.Ctx) error) {
	handler := httpImageApiDelivery{
		imageService: imageService,
	}
	group := app.Group("images", middlewares...)

	group.Post("/", handler.UploadImage)
}

func (i *httpImageApiDelivery) UploadImage(ctx *fiber.Ctx) error {
	wg := ctx.Locals("wg").(*sync.WaitGroup)
	wg.Add(1)
	defer wg.Done()
	userId := middleware.GetUserId(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parse form", err)
	}
	if file == nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "No image detected", nil)
	}
	fileContents, err := file.Open()
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to open file", err)
	}
	defer fileContents.Close()

	buffer := make([]byte, 512)
	_, err = fileContents.Read(buffer)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to read buffer", err)
	}
	fileType := mimetype.Detect(buffer)
	if filepath.Ext(file.Filename) == ".gif" {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "file type must be image", err)
	}

	if !strings.HasPrefix(fileType.String(), "image/") {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "file type must be image", err)
	}
	res, err := i.imageService.UploadImage(file, userId)
	if err != nil {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to upload file", err)
	}
	return fiber_response.ReturnStatusCreated(ctx, "Success", res)
}
