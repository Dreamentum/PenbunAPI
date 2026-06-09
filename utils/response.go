package utils

import (
	"github.com/gofiber/fiber/v2"
	"PenbunAPI/models"
)

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func FailResponse(c *fiber.Ctx, message string) error {
	return c.Status(400).JSON(models.ApiResponse{
		Status:  "fail",
		Message: message,
	})
}

func ErrorResponse(c *fiber.Ctx, message string) error {
	return c.Status(500).JSON(models.ApiResponse{
		Status:  "error",
		Message: message,
	})
}

func LoginSuccessResponse(c *fiber.Ctx, token, message string) error {
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Token:   token,
		Message: message,
	})
}

func LoginFailResponse(c *fiber.Ctx, message string) error {
	return c.Status(401).JSON(models.ApiResponse{
		Status:  "fail",
		Message: message,
	})
}

func UnauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(401).JSON(models.ApiResponse{
		Status:  "error",
		Message: "Unauthorized",
	})
}
