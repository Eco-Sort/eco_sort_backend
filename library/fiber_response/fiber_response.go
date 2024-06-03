package fiber_response

import "github.com/gofiber/fiber/v2"

func ReturnStatusUnprocessableEntity(c *fiber.Ctx, messages string, errorData any) error {
	statusCode := fiber.StatusUnprocessableEntity
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": messages,
	})
}

func ReturnStatusOk(c *fiber.Ctx, messages string, data any) error {
	statusCode := fiber.StatusOK
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": messages,
		"data":    data,
	})
}

func ReturnStatusUnauthorized(c *fiber.Ctx) error {
	statusCode := fiber.StatusUnauthorized
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": "Unauthorized",
	})
}

func ReturnStatusServerError(c *fiber.Ctx, messages string, err error) error {
	statusCode := fiber.StatusInternalServerError
	return c.Status(statusCode).JSON(fiber.Map{
		"status":  statusCode,
		"message": messages,
	})
}

func ReturnStatusCreated(c *fiber.Ctx, messages string, data any) error {
	statusCode := fiber.StatusCreated
	return c.Status(statusCode).JSON(fiber.Map{
		"status":   statusCode,
		"messages": messages,
		"data":     data,
	})
}
