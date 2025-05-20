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
	// Group สำหรับ Book API [ver 1.5.9]
	book := protected.Group("/book")
	book.Post("/insert", controllers.InsertBook)                 // เพิ่ม Book
	book.Get("/select/all", controllers.SelectAllBooks)          // ดึงข้อมูล Book ทั้งหมด (ไม่มี Paging)
	book.Get("/select/page", controllers.SelectPageBooks)        // ดึงข้อมูล Book แบบ Paging
	book.Get("/select/:id", controllers.SelectBookByID)          // ดึงข้อมูล Book ตาม book_code
	book.Get("/select/name/:name", controllers.SelectBookByName) // ดึงข้อมูล Book ตามชื่อหนังสือแบบ LIKE
	book.Put("/update/:id", controllers.UpdateBookByID)          // อัปเดต Book ตาม book_code
	book.Put("/delete/:id", controllers.DeleteBookByID)          // Soft Delete Book (is_delete = 1)
	book.Delete("/remove/:id", controllers.RemoveBookByID)       // ลบข้อมูลจริง (Hard Delete)

	// Group สำหรับ Publisher Type API [ver 1.5.1]
	publisherType := protected.Group("/publishertype")
	publisherType.Post("/insert", controllers.InsertPublisherType)            // เพิ่ม Publisher Type
	publisherType.Get("/select/all", controllers.SelectAllPublisherTypes)     // ดึงข้อมูล Publisher Type ทั้งหมด
	publisherType.Get("/select/page", controllers.SelectPagePublisherTypes)   // ดึงข้อมูล Publisher Type แบบ Paging
	publisherType.Get("/select/:id", controllers.SelectPublisherTypeByID)     // ดึงข้อมูล Publisher Type ตาม ID
	publisherType.Get("/select/:name", controllers.SelectPublisherTypeByName) // ดึงข้อมูล Publisher Type ตาม Name [version 1.6.3]
	publisherType.Put("/update/:id", controllers.UpdatePublisherTypeByID)     // อัปเดต Publisher Type ตาม ID
	publisherType.Put("/delete/:id", controllers.DeletePublisherTypeByID)     // เปลี่ยน is_delete = 1
	publisherType.Delete("/remove/:id", controllers.RemovePublisherTypeByID)  // ลบข้อมูลจริง

	// Group สำหรับ Publisher API [ver 1.5.1]
	publisher := protected.Group("/publisher")
	publisher.Post("/insert", controllers.InsertPublisher)            // เพิ่ม Publisher
	publisher.Get("/select/all", controllers.SelectAllPublisher)      // ดึงข้อมูล Publisher ทั้งหมด (ไม่มี Paging)
	publisher.Get("/select/page", controllers.SelectPagePublisher)    // ดึงข้อมูล Publisher ทั้งหมด (รองรับ Paging)
	publisher.Get("/select/:id", controllers.SelectPublisherByID)     // ดึงข้อมูล Publisher ตาม ID
	publisher.Get("/select/:name", controllers.SelectPublisherByName) // ดึงข้อมูล Publisher ตาม Name [version 1.6.3]
	publisher.Put("/update/:id", controllers.UpdatePublisherByID)     // อัปเดต Publisher ตาม ID
	publisher.Put("/delete/:id", controllers.DeletePublisherByID)     // เปลี่ยน is_delete = 1
	publisher.Delete("/remove/:id", controllers.RemovePublisherByID)  // ลบข้อมูลจริง

	// Group สำหรับ Customer Type API [ver 1.5.3]
	customerType := protected.Group("/customertype")
	customerType.Post("/insert", controllers.InsertCustomerType)
	customerType.Get("/select/all", controllers.SelectAllCustomerTypes)
	customerType.Get("/select/page", controllers.SelectPageCustomerTypes)
	customerType.Get("/select/:id", controllers.SelectCustomerTypeByID)
	customerType.Put("/update/:id", controllers.UpdateCustomerTypeByID)
	customerType.Put("/delete/:id", controllers.DeleteCustomerTypeByID)
	customerType.Delete("/remove/:id", controllers.RemoveCustomerTypeByID)

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
	bookType := protected.Group("/booktype")
	bookType.Post("/insert", controllers.InsertBookType)
	bookType.Get("/select/all", controllers.SelectAllBookTypes)
	bookType.Get("/select/page", controllers.SelectPageBookTypes)
	bookType.Get("/select/:id", controllers.SelectBookTypeByID)
	bookType.Put("/update/:id", controllers.UpdateBookTypeByID)
	bookType.Put("/delete/:id", controllers.DeleteBookTypeByID)
	bookType.Delete("/remove/:id", controllers.RemoveBookTypeByID)

	// Group สำหรับ Discount Type API [ver 1.5.5]
	discountType := protected.Group("/discounttype")
	discountType.Post("/insert", controllers.InsertDiscountType)
	discountType.Get("/select/all", controllers.SelectAllDiscountType)
	discountType.Get("/select/page", controllers.SelectPageDiscountType)
	discountType.Get("/select/:id", controllers.SelectDiscountTypeByID)
	discountType.Put("/update/:id", controllers.UpdateDiscountTypeByID)
	discountType.Put("/delete/:id", controllers.DeleteDiscountTypeByID)
	discountType.Delete("/remove/:id", controllers.RemoveDiscountTypeByID)

	// Group สำหรับ Discount API [ver 1.5.7]
	discount := protected.Group("/discount")
	discount.Post("/insert", controllers.InsertDiscount)
	discount.Get("/select/all", controllers.SelectAllDiscount)
	discount.Get("/select/page", controllers.SelectPageDiscount)
	discount.Get("/select/:id", controllers.SelectDiscountByID)
	discount.Put("/update/:id", controllers.UpdateDiscountByID)
	discount.Put("/delete/:id", controllers.DeleteDiscountByID)
	discount.Delete("/remove/:id", controllers.RemoveDiscountByID)

}
