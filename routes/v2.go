package routes

import (
	"os"
	"database/sql"
	"PenbunAPI/controllers"
	"PenbunAPI/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterV2Routes(app *fiber.App, db *sql.DB) {
	// Group สำหรับ API Version 2
	v2 := app.Group("/api/v2")

	// Group สำหรับ public API
	public := v2.Group("/public")
	RegisterPublicRoutes(public, db)

	// Group สำหรับ protected API
	protected := v2.Group("/protected")

	protected.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	protected.Get("/books", controllers.GetAllBooks) // ดึงหนังสือทั้งหมด
	protected.Get("/books/:id", controllers.GetBook) // ดึงหนังสือตาม ID
}
