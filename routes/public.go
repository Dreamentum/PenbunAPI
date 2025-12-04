package routes

import (
	"database/sql"
	"log"
	"github.com/gofiber/fiber/v2"
)



// RegisterPublicRoutes จะลงทะเบียน Route สำหรับ public API
func RegisterPublicRoutes(group fiber.Router, db *sql.DB) {
	// Route "/"
	group.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the PENBUN API v1",
			"status":  "success",
		})
	})

	// Hello Route handles the /hello endpoint
	group.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	// Route "/welcome"
	group.Get("/welcome", func(c *fiber.Ctx) error {
		var version string
		err := db.QueryRow("SELECT @@VERSION").Scan(&version)
		if err != nil {
			log.Println("[ERROR] Failed to query database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve SQL Server version",
			})
		}
		return c.JSON(fiber.Map{
			"sql_version": version,
		})
	})

}

