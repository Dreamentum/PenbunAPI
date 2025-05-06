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

	// Group สำหรับ Public API [ver 1.0.1]
	public := v1.Group("/public")
	RegisterPublicRoutes(public, db)

	// Group สำหรับ login/logout API [ver 1.0.1]
	public.Post("/login", controllers.Login)   // Route สำหรับ login (ไม่ใช้ Middleware)
	public.Post("/logout", controllers.Logout) // Route สำหรับ logout (ไม่ใช้ Middleware)

	// Group สำหรับ Protected API [ver 1.0.1]
	protected := v1.Group("/protected")
	protected.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	protected.Post("/refresh", controllers.RefreshToken)  // Route สำหรับ Refresh Token
	protected.Get("/reference", controllers.GetReference) // Route สำหรับ get ค่า references

	// Routes สำหรับการจัดการ books [ver 1.0.1]
	protected.Post("/books", controllers.CreateBook)
	protected.Put("/books/:id", controllers.UpdateBook)
	protected.Delete("/books/:id", controllers.DeleteBook)

	// Group สำหรับ Publisher API [ver 1.5.1]
	publishers := protected.Group("/publisher")
	publishers.Post("/insert", controllers.InsertPublisher)           // เพิ่ม Publisher
	publishers.Get("/select/all", controllers.SelectAllPublisher)     // ดึงข้อมูล Publisher ทั้งหมด (ไม่มี Paging)
	publishers.Get("/select/page", controllers.SelectPagePublisher)   // ดึงข้อมูล Publisher ทั้งหมด (รองรับ Paging)
	publishers.Get("/select/:id", controllers.SelectPublisherByID)    // ดึงข้อมูล Publisher ตาม ID
	publishers.Put("/update/:id", controllers.UpdatePublisherByID)    // อัปเดต Publisher ตาม ID
	publishers.Put("/delete/:id", controllers.DeletePublisherByID)    // เปลี่ยน is_delete = 1
	publishers.Delete("/remove/:id", controllers.RemovePublisherByID) // ลบข้อมูลจริง

	// Group สำหรับ Publisher Type API [ver 1.5.1]
	publisherTypes := protected.Group("/publishertype")
	publisherTypes.Post("/insert", controllers.InsertPublisherType)           // เพิ่ม Publisher Type
	publisherTypes.Get("/select/all", controllers.SelectAllPublisherTypes)    // ดึงข้อมูล Publisher Type ทั้งหมด
	publisherTypes.Get("/select/page", controllers.SelectPagePublisherTypes)  // ดึงข้อมูล Publisher Type แบบ Paging
	publisherTypes.Get("/select/:id", controllers.SelectPublisherTypeByID)    // ดึงข้อมูล Publisher Type ตาม ID
	publisherTypes.Put("/update/:id", controllers.UpdatePublisherTypeByID)    // อัปเดต Publisher Type ตาม ID
	publisherTypes.Put("/delete/:id", controllers.DeletePublisherTypeByID)    // เปลี่ยน is_delete = 1
	publisherTypes.Delete("/remove/:id", controllers.RemovePublisherTypeByID) // ลบข้อมูลจริง

	// Group สำหรับ Customer Type API [ver 1.5.3]
	customerTypes := protected.Group("/customertype")
	customerTypes.Post("/insert", controllers.InsertCustomerType)
	customerTypes.Get("/select/all", controllers.SelectAllCustomerTypes)
	customerTypes.Get("/select/page", controllers.SelectPageCustomerTypes)
	customerTypes.Get("/select/:id", controllers.SelectCustomerTypeByID)
	customerTypes.Put("/update/:id", controllers.UpdateCustomerTypeByID)
	customerTypes.Put("/delete/:id", controllers.DeleteCustomerTypeByID)
	customerTypes.Delete("/remove/:id", controllers.RemoveCustomerTypeByID)

	// Group สำหรับ Customer API [ver 1.5.3]
	customer := protected.Group("/customer")
	customer.Post("/insert", controllers.InsertCustomer)           // เพิ่ม Customer
	customer.Get("/select/all", controllers.SelectAllCustomers)    // ดึงข้อมูล Customer ทั้งหมด (ไม่มี Paging)
	customer.Get("/select/page", controllers.SelectPageCustomers)  // ดึงข้อมูล Customer ทั้งหมด (รองรับ Paging)
	customer.Get("/select/:id", controllers.SelectCustomerByID)    // ดึงข้อมูล Customer ตาม ID
	customer.Put("/update/:id", controllers.UpdateCustomerByID)    // อัปเดต Customer ตาม ID
	customer.Put("/delete/:id", controllers.DeleteCustomerByID)    // เปลี่ยน is_delete = 1
	customer.Delete("/remove/:id", controllers.RemoveCustomerByID) // ลบข้อมูลจริง

	// Group สำหรับ Book Type API [ver 1.5.3]
	bookTypes := protected.Group("/booktype")
	bookTypes.Post("/insert", controllers.InsertBookType)
	bookTypes.Get("/select/all", controllers.SelectAllBookTypes)
	bookTypes.Get("/select/page", controllers.SelectPageBookTypes)
	bookTypes.Get("/select/:id", controllers.SelectBookTypeByID)
	bookTypes.Put("/update/:id", controllers.UpdateBookTypeByID)
	bookTypes.Put("/delete/:id", controllers.DeleteBookTypeByID)
	bookTypes.Delete("/remove/:id", controllers.RemoveBookTypeByID)

	// Group สำหรับ Discount Type API [ver 1.5.5]
	discountTypes := protected.Group("/discounttype")
	discountTypes.Post("/insert", controllers.InsertDiscountType)
	discountTypes.Get("/select/all", controllers.SelectAllDiscountType)
	discountTypes.Get("/select/page", controllers.SelectPageDiscountType)
	discountTypes.Get("/select/:id", controllers.SelectDiscountTypeByID)
	discountTypes.Put("/update/:id", controllers.UpdateDiscountTypeByID)
	discountTypes.Put("/delete/:id", controllers.DeleteDiscountTypeByID)
	discountTypes.Delete("/remove/:id", controllers.RemoveDiscountTypeByID)

}
