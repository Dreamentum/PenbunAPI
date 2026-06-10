package middleware

import (
	"github.com/gofiber/fiber/v2"

	"PenbunAPI/models"
)

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(models.ApiResponse{
		Status:  "error",
		Message: err.Error(),
	})
}
