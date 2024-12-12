package routes

import (
	"PenbunAPI/controllers"
	"PenbunAPI/middleware"
	"os"

	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// RegisterV1Routes will register all V1 routes
func RegisterV1Routes(app *fiber.App, db *sql.DB) {
	// Group สำหรับ API Version 1
	v1 := app.Group("/api/v1")

	// Group สำหรับ Public API
	public := v1.Group("/public")
	RegisterPublicRoutes(public, db)

	public.Post("/login", controllers.Login)   // Route สำหรับ login (ไม่ใช้ Middleware)
	public.Post("/logout", controllers.Logout) // Route สำหรับ logout (ไม่ใช้ Middleware)

	// Group สำหรับ Protected API
	protected := v1.Group("/protected")
	protected.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	protected.Post("/refresh", controllers.RefreshToken)  // Route สำหรับ Refresh Token
	protected.Get("/reference", controllers.GetReference) // Route สำหรับ get ค่า references

	// Routes สำหรับการจัดการ books
	protected.Post("/books", controllers.CreateBook)
	protected.Put("/books/:id", controllers.UpdateBook)
	protected.Delete("/books/:id", controllers.DeleteBook)

	// Group สำหรับ Publisher API
	publishers := protected.Group("/publishers")
	publishers.Post("/insert", controllers.InsertPublisher)           // เพิ่ม Publisher
	publishers.Get("/select/all", controllers.SelectAllPublishers)    // ดึงข้อมูล Publisher ทั้งหมด
	publishers.Get("/select/:id", controllers.SelectPublisherByID)    // ดึงข้อมูล Publisher ตาม ID
	publishers.Put("/update/:id", controllers.UpdatePublisherByID)    // อัปเดต Publisher ตาม ID
	publishers.Put("/delete/:id", controllers.DeletePublisherByID)    // เปลี่ยน is_delete = 1
	publishers.Delete("/remove/:id", controllers.RemovePublisherByID) // ลบข้อมูลจริง
}
