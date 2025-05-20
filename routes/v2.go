package routes

import (
	"PenbunAPI/controllers"
	"PenbunAPI/middleware"
	"database/sql"
	"os"

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
	protected.Post("/refresh", controllers.RefreshToken) // Route สำหรับ Refresh Token
}
