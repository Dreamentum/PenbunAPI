package utils

import "github.com/gofiber/fiber/v2"

// ResolveUser extracts the username from query parameters or returns "UNKNOWN"
func ResolveUser(c *fiber.Ctx) string {
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}
	return username
}
