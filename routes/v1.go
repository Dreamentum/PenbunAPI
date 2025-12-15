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
	public.Post("/login", controllers.Login) // Route สำหรับ login (ไม่ใช้ Middleware)
	// Apply Middleware to logout to enable user logging
	public.Post("/logout", middleware.JWTMiddleware(os.Getenv("JWT_SECRET")), controllers.Logout)

	// Group สำหรับ Protected API [ver 1.0.1]
	protected := v1.Group("/protected")
	protected.Use(middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	protected.Post("/refresh", controllers.RefreshToken)  // Route สำหรับ Refresh Token
	protected.Get("/reference", controllers.GetReference) // Route สำหรับ get ค่า references

	// Group สำหรับ Vendor API [ver 2.3.0]
	vendor := protected.Group("/vendor")
	vendor.Post("/insert", controllers.InsertVendor)
	vendor.Get("/select/all", controllers.SelectAllVendors)
	vendor.Get("/select/page", controllers.SelectPageVendors)
	vendor.Get("/select/:id", controllers.SelectVendorByID)
	vendor.Get("/select/name/:name", controllers.SelectVendorByName)
	vendor.Put("/update/:id", controllers.UpdateVendorByID)
	vendor.Put("/delete/:id", controllers.DeleteVendorByID)
	vendor.Delete("/remove/:id", controllers.RemoveVendorByID)

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

	// Group สำหรับ Discount Type API [ver 1.5.5]
	discountType := protected.Group("/discounttype")
	discountType.Post("/insert", controllers.InsertDiscountType)
	discountType.Get("/select/all", controllers.SelectAllDiscountType)
	discountType.Get("/select/page", controllers.SelectPageDiscountType)
	discountType.Get("/select/:id", controllers.SelectDiscountTypeByID)
	discountType.Get("/select/name/:name", controllers.SelectDiscountTypeByName)
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

	// Group สำหรับ Vendor Type API [ver 1.7.1]
	vendorType := protected.Group("/vendortype")
	vendorType.Post("/insert", controllers.InsertVendorType)                 // เพิ่ม Vendor Type
	vendorType.Get("/select/all", controllers.SelectAllVendorType)           // ดึงข้อมูล Vendor Type ทั้งหมด
	vendorType.Get("/select/page", controllers.SelectPageVendorType)         // ดึงข้อมูล Vendor Type แบบ Paging
	vendorType.Get("/select/:id", controllers.SelectVendorTypeByID)          // ดึงข้อมูล Vendor Type ตาม ID
	vendorType.Get("/select/name/:name", controllers.SelectVendorTypeByName) // ดึงข้อมูล Vendor Type ตาม Name
	vendorType.Put("/update/:id", controllers.UpdateVendorTypeByID)          // อัปเดต Vendor Type ตาม ID
	vendorType.Put("/delete/:id", controllers.DeleteVendorTypeByID)          // เปลี่ยน is_delete = 1
	vendorType.Delete("/remove/:id", controllers.RemoveVendorTypeByID)       // ลบข้อมูลจริง

	// Group สำหรับ Vendor API [ver 2.3.0]
	vendorGroup := protected.Group("/vendor")
	vendorGroup.Post("/insert", controllers.InsertVendor)
	vendorGroup.Get("/select/all", controllers.SelectAllVendors)
	vendorGroup.Get("/select/page", controllers.SelectPageVendors)
	vendorGroup.Get("/select/:id", controllers.SelectVendorByID)
	vendorGroup.Get("/select/name/:name", controllers.SelectVendorByName)
	vendorGroup.Put("/update/:id", controllers.UpdateVendorByID)
	vendorGroup.Put("/delete/:id", controllers.DeleteVendorByID)
	vendorGroup.Delete("/remove/:id", controllers.RemoveVendorByID)

	// Group สำหรับ Unit Type API [ver 1.7.3]
	unitType := protected.Group("/unittype")
	unitType.Post("/insert", controllers.InsertUnitType)                 // เพิ่ม Unit Type
	unitType.Get("/select/all", controllers.SelectAllUnitType)           // ดึงข้อมูลทั้งหมด
	unitType.Get("/select/page", controllers.SelectPageUnitType)         // ดึงข้อมูลแบบ Paging
	unitType.Get("/select/:id", controllers.SelectUnitTypeByID)          // ดึงข้อมูลตาม ID
	unitType.Get("/select/name/:name", controllers.SelectUnitTypeByName) // ดึงข้อมูลตาม Name
	unitType.Put("/update/:id", controllers.UpdateUnitTypeByID)          // อัปเดตตาม ID
	unitType.Put("/delete/:id", controllers.DeleteUnitTypeByID)          // Soft Delete
	unitType.Delete("/remove/:id", controllers.RemoveUnitTypeByID)       // ลบข้อมูลจริง

	// Group สำหรับ Product Group API [ver 2.3.0]
	productGroup := protected.Group("/productgroup")
	productGroup.Post("/insert", controllers.InsertProductGroup)
	productGroup.Get("/select/all", controllers.SelectAllProductGroup)
	productGroup.Get("/select/page", controllers.SelectPageProductGroup)
	productGroup.Get("/select/:id", controllers.SelectProductGroupByID)
	productGroup.Get("/select/name/:name", controllers.SelectProductGroupByName)
	productGroup.Put("/update/:id", controllers.UpdateProductGroupByID)
	productGroup.Put("/delete/:id", controllers.DeleteProductGroupByID)
	productGroup.Delete("/remove/:id", controllers.RemoveProductGroupByID)

	// Group สำหรับ Product Category API [ver 1.8.3]
	productCategory := protected.Group("/productcategory")
	productCategory.Post("/insert", controllers.InsertProductCategory)
	productCategory.Get("/select/all", controllers.SelectAllProductCategory)
	productCategory.Get("/select/page", controllers.SelectPageProductCategory)
	productCategory.Get("/select/:id", controllers.SelectProductCategoryByID)
	productCategory.Get("/select/name/:name", controllers.SelectProductCategoryByName)
	productCategory.Put("/update/:id", controllers.UpdateProductCategoryByID)
	productCategory.Put("/delete/:id", controllers.DeleteProductCategoryByID)
	productCategory.Delete("/remove/:id", controllers.RemoveProductCategoryByID)

	// Group สำหรับ Product Format Type API [ver 1.8.1]
	productFormatType := protected.Group("/productformattype")
	productFormatType.Post("/insert", controllers.InsertProductFormatType)                 // เพิ่ม Product Format Type
	productFormatType.Get("/select/all", controllers.SelectAllProductFormatType)           // ดึงข้อมูลทั้งหมด
	productFormatType.Get("/select/page", controllers.SelectPageProductFormatType)         // ดึงข้อมูลแบบ Paging
	productFormatType.Get("/select/:id", controllers.SelectProductFormatTypeByID)          // ดึงข้อมูลตาม ID
	productFormatType.Get("/select/name/:name", controllers.SelectProductFormatTypeByName) // ดึงข้อมูลตาม Name
	productFormatType.Put("/update/:id", controllers.UpdateProductFormatTypeByID)          // อัปเดตตาม ID
	productFormatType.Put("/delete/:id", controllers.DeleteProductFormatTypeByID)          // Soft Delete
	productFormatType.Delete("/remove/:id", controllers.RemoveProductFormatTypeByID)       // ลบข้อมูลจริง

	// Group สำหรับ Product Pack Config API [ver 1.8.2]
	productPackConfig := protected.Group("/productpackconfig")
	productPackConfig.Post("/insert", controllers.InsertProductPackConfig)                 // เพิ่ม Product Pack Config
	productPackConfig.Get("/select/all", controllers.SelectAllProductPackConfig)           // ดึงข้อมูลทั้งหมด
	productPackConfig.Get("/select/page", controllers.SelectPageProductPackConfig)         // ดึงข้อมูลแบบ Paging
	productPackConfig.Get("/select/:id", controllers.SelectProductPackConfigByID)          // ดึงข้อมูลตาม ID
	productPackConfig.Get("/select/name/:name", controllers.SelectProductPackConfigByName) // ดึงข้อมูลตาม Name (Note/ProductID)
	productPackConfig.Put("/update/:id", controllers.UpdateProductPackConfigByID)          // อัปเดตตาม ID
	productPackConfig.Put("/delete/:id", controllers.DeleteProductPackConfigByID)          // Soft Delete
	productPackConfig.Delete("/remove/:id", controllers.RemoveProductPackConfigByID)       // ลบข้อมูลจริง

	// Group สำหรับ Product API [ver 2.0.0]
	product := protected.Group("/product")
	product.Post("/insert", controllers.InsertProduct)                 // เพิ่ม Product (with dummy ID 'TEMP')
	product.Get("/select/all", controllers.SelectAllProducts)          // ดึงข้อมูล Product ทั้งหมด
	product.Get("/select/page", controllers.SelectPageProducts)        // ดึงข้อมูล Product แบบ Paging
	product.Get("/select/:id", controllers.SelectProductByID)          // ดึงข้อมูล Product ตาม ID
	product.Get("/select/name/:name", controllers.SelectProductByName) // ดึงข้อมูล Product ตาม Name (TH or EN)
	product.Put("/update/:id", controllers.UpdateProductByID)          // อัปเดต Product ตาม ID
	product.Put("/delete/:id", controllers.DeleteProductByID)          // Soft Delete
	product.Delete("/remove/:id", controllers.RemoveProductByID)       // ลบข้อมูลจริง

	// Group สำหรับ Warehouse API [ver 2.3.0]
	warehouse := protected.Group("/warehouse")
	warehouse.Post("/insert", controllers.InsertWarehouse)
	warehouse.Get("/select/all", controllers.SelectAllWarehouse)
	warehouse.Get("/select/page", controllers.SelectPageWarehouse)
	warehouse.Get("/select/:id", controllers.SelectWarehouseByID)
	warehouse.Get("/select/name/:name", controllers.SelectWarehouseByName)
	warehouse.Put("/update/:id", controllers.UpdateWarehouseByID)
	warehouse.Put("/delete/:id", controllers.DeleteWarehouseByID)
	warehouse.Delete("/remove/:id", controllers.RemoveWarehouseByID)

	// Group สำหรับ Receive API [ver 2.3.0]
	receive := protected.Group("/receive")
	receive.Post("/insert", controllers.InsertReceiveNote)
	receive.Get("/select/all", controllers.SelectAllReceiveNotes)
	receive.Get("/select/page", controllers.SelectPageReceiveNotes)
	receive.Get("/select/:id", controllers.SelectReceiveNoteByID)
	receive.Put("/update/:id", controllers.UpdateReceiveNoteByID)
	receive.Put("/delete/:id", controllers.DeleteReceiveNoteByID)
	receive.Delete("/remove/:id", controllers.RemoveReceiveNoteByID)

	// Group สำหรับ Order API [ver 2.3.0]
	order := protected.Group("/order")
	order.Post("/insert", controllers.InsertOrder)
	order.Get("/select/all", controllers.SelectAllOrders)
	order.Get("/select/page", controllers.SelectPageOrders)
	order.Get("/select/:id", controllers.SelectOrderByID)
	order.Put("/update/:id", controllers.UpdateOrderByID)
	order.Put("/delete/:id", controllers.DeleteOrderByID)
	order.Delete("/remove/:id", controllers.RemoveOrderByID)

}
