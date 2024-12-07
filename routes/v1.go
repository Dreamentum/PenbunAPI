package routes

import (
	"PenbunAPI/controllers"
	"PenbunAPI/middleware"
	"os"

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

	RegisterPublisherRoutes(v1)
}

func RegisterPublisherRoutes(v1 fiber.Router) {
	publisher := v1.Group("/publisher")
	publisher.Post("/insert", controllers.InsertPublisher)
	publisher.Get("/select/all", controllers.SelectAllPublishers)
	publisher.Get("/select/:id", controllers.SelectPublisherByID)
	publisher.Put("/update/:id", controllers.UpdatePublisherByID)
	publisher.Put("/delete/:id", controllers.DeletePublisher)
	publisher.Delete("/remove/:id", controllers.RemovePublisher)
}
