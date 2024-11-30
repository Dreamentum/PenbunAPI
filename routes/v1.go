package routes

import (
	"os"
	"PenbunAPI/middleware"
	"PenbunAPI/controllers"

	"database/sql"
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes will register all V1 routes
func RegisterV1Routes(app *fiber.App, db *sql.DB) {
	// Group สำหรับ API Version 1
	v1 := app.Group("/api/v1")

	// Group สำหรับ Public API
	public := v1.Group("/public")
	RegisterPublicRoutes(public, db)

	// Route สำหรับ login (ไม่ใช้ Middleware)
	public.Post("/login", controllers.Login)

	// Route สำหรับ logout (ไม่ใช้ Middleware)
	public.Post("/logout", controllers.Logout)

	// Group สำหรับ protected API
	protected := v1.Group("/protected")
	protected.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	// Route สำหรับ Refresh Token
	protected.Post("/refresh", controllers.RefreshToken)

	// Route สำหรับ get ค่า references
	protected.Get("/reference", controllers.GetReference)

	// Routes สำหรับการจัดการ books
	protected.Post("/books", controllers.CreateBook)
	protected.Put("/books/:id", controllers.UpdateBook)
	protected.Delete("/books/:id", controllers.DeleteBook)
}
