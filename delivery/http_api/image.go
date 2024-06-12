package http_api

import (
	"fmt"
	"os"
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
	imageService   domain.ImageService
	garbageService domain.GarbageService
}

func NewImageHttpApiDelivery(app fiber.Router, imageService domain.ImageService, garbageService domain.GarbageService, middlewares ...func(*fiber.Ctx) error) {
	handler := httpImageApiDelivery{
		imageService:   imageService,
		garbageService: garbageService,
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
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to parse form", err)
	}
	if file == nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "No image detected", nil)
	}
	fileContents, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to open file", err)
	}
	defer fileContents.Close()

	buffer := make([]byte, 512)
	_, err = fileContents.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to read buffer", err)
	}
	fileType := mimetype.Detect(buffer)
	if filepath.Ext(file.Filename) == ".gif" {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "file type must be image", nil)
	}

	if !strings.HasPrefix(fileType.String(), "image/") {
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "file type must be image", nil)
	}

	res, err := i.imageService.ProcessImage(file)
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to upload file", err)
	}
	filePath := domain.MLImage
	fileLocal, err := os.Open(filePath + "/" + res.Filename)
	if err != nil {
		panic(err)
	}
	defer fileLocal.Close()
	resp, err := i.imageService.UploadImage(fileLocal, filepath.Base(filePath), userId)
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to upload file", err)
	}

	garbage := domain.Garbage{
		UserID: userId,
		Image:  *resp,
	}
	var imageObject = []domain.ImageObject{}
	for _, s := range res.Result {
		imageObject = append(imageObject, domain.ImageObject{
			CategoryID: uint(s.Class),
			Confidence: s.Score,
		})
	}
	garbage.ImageObject = append(garbage.ImageObject, imageObject...)

	resGarbage, err := i.garbageService.Create(garbage)
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to create garbage", err)
	}
	getGarbage, err := i.garbageService.GetById(resGarbage)
	if err != nil {
		fmt.Println(err)
		return fiber_response.ReturnStatusUnprocessableEntity(ctx, "Failed to getting garbage", err)
	}
	var imageObjects []map[string]any
	for _, imageObject := range garbage.ImageObject {
		imageObjects = append(imageObjects, map[string]any{
			"classification": imageObject.Category,
			"sorting":        imageObject.Category.Sorting,
			"confidence":     imageObject.Confidence,
		})
	}
	statusCode := fiber.StatusOK
	return ctx.Status(statusCode).JSON(fiber.Map{
		"error":        false,
		"messages":     "Success",
		"image":        getGarbage.Image,
		"image_detail": imageObjects,
	})
}
