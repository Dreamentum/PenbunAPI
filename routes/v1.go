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
	publishers.Get("/select/all", controllers.SelectAllPublishers)    // ดึงข้อมูล Publisher ทั้งหมด (ไม่มี Paging)
	publishers.Get("/select/page", controllers.SelectPagePublishers)  // ดึงข้อมูล Publisher ทั้งหมด (รองรับ Paging)
	publishers.Get("/select/:id", controllers.SelectPublisherByID)    // ดึงข้อมูล Publisher ตาม ID
	publishers.Put("/update/:id", controllers.UpdatePublisherByID)    // อัปเดต Publisher ตาม ID
	publishers.Put("/delete/:id", controllers.DeletePublisherByID)    // เปลี่ยน is_delete = 1
	publishers.Delete("/remove/:id", controllers.RemovePublisherByID) // ลบข้อมูลจริง

	// Group สำหรับ Publisher Type API
	publisherTypes := protected.Group("/publishertype")
	publisherTypes.Post("/insert", controllers.InsertPublisherType)           // เพิ่ม Publisher Type
	publisherTypes.Get("/select/all", controllers.SelectAllPublisherTypes)    // ดึงข้อมูล Publisher Type ทั้งหมด
	publisherTypes.Get("/select/page", controllers.SelectPagePublisherTypes)  // ดึงข้อมูล Publisher Type แบบ Paging
	publisherTypes.Get("/select/:id", controllers.SelectPublisherTypeByID)    // ดึงข้อมูล Publisher Type ตาม ID
	publisherTypes.Put("/update/:id", controllers.UpdatePublisherTypeByID)    // อัปเดต Publisher Type ตาม ID
	publisherTypes.Put("/delete/:id", controllers.DeletePublisherTypeByID)    // เปลี่ยน is_delete = 1
	publisherTypes.Delete("/remove/:id", controllers.RemovePublisherTypeByID) // ลบข้อมูลจริง

	// Group สำหรับ Customer Type API
	customerTypes := protected.Group("/customertype")
	customerTypes.Post("/insert", controllers.InsertCustomerType)
	customerTypes.Get("/select/all", controllers.SelectAllCustomerTypes)
	customerTypes.Get("/select/page", controllers.SelectPageCustomerTypes)
	customerTypes.Get("/select/:id", controllers.SelectCustomerTypeByID)
	customerTypes.Put("/update/:id", controllers.UpdateCustomerTypeByID)
	customerTypes.Put("/delete/:id", controllers.DeleteCustomerTypeByID)
	customerTypes.Delete("/remove/:id", controllers.RemoveCustomerTypeByID)

}
