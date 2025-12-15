package controllers

import "github.com/gofiber/fiber/v2"

// resolveUser extracts the username from query parameters or returns "UNKNOWN"
func resolveUser(c *fiber.Ctx) string {
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}
	return username
}
