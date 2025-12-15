package controllers

import (
	"PenbunAPI/models"
	"database/sql"
	
	"github.com/gofiber/fiber/v2"
)

// SelectAllOrders ดึงข้อมูลใบสั่งขายทั้งหมด
func SelectAllOrders(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)

	query := `
		SELECT o.order_id, o.customer_id, c.customer_name, o.warehouse_id, w.warehouse_name,
		       o.doc_date, o.doc_type, o.total_amount, o.discount_amount, o.net_amount, o.vat_amount, o.grand_total,
		       o.update_by, o.update_date, o.is_active
		FROM tb_order o
		LEFT JOIN tb_customer c ON o.customer_id = c.customer_id
		LEFT JOIN tb_warehouse w ON o.warehouse_id = w.warehouse_id
		WHERE o.is_delete = 0
		ORDER BY o.doc_date DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var customerName, warehouseName *string
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &customerName, &o.WarehouseID, &warehouseName,
			&o.DocDate, &o.DocType, &o.TotalAmount, &o.DiscountAmount, &o.NetAmount, &o.VatAmount, &o.GrandTotal,
			&o.UpdateBy, &o.UpdateDate, &o.IsActive); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		o.CustomerName = customerName
		orders = append(orders, o)
	}

	return c.JSON(fiber.Map{"data": orders})
}

// SelectPageOrders ดึงข้อมูลใบสั่งขายแบบ Paging
func SelectPageOrders(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM tb_order WHERE is_delete = 0`
	if err := db.QueryRow(countQuery).Scan(&total); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	query := `
		SELECT o.order_id, o.customer_id, c.customer_name, o.warehouse_id, w.warehouse_name,
		       o.doc_date, o.doc_type, o.total_amount, o.discount_amount, o.net_amount, o.vat_amount, o.grand_total,
		       o.update_by, o.update_date, o.is_active
		FROM tb_order o
		LEFT JOIN tb_customer c ON o.customer_id = c.customer_id
		LEFT JOIN tb_warehouse w ON o.warehouse_id = w.warehouse_id
		WHERE o.is_delete = 0
		ORDER BY o.doc_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var customerName, warehouseName *string
		if err := rows.Scan(&o.OrderID, &o.CustomerID, &customerName, &o.WarehouseID, &warehouseName,
			&o.DocDate, &o.DocType, &o.TotalAmount, &o.DiscountAmount, &o.NetAmount, &o.VatAmount, &o.GrandTotal,
			&o.UpdateBy, &o.UpdateDate, &o.IsActive); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		o.CustomerName = customerName
		orders = append(orders, o)
	}

	return c.JSON(fiber.Map{
		"data":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// SelectOrderByID ดึงข้อมูลใบสั่งขายตาม ID
func SelectOrderByID(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)
	id := c.Params("id")

	// 1. Header
	queryHeader := `
		SELECT o.order_id, o.customer_id, c.customer_name, o.warehouse_id, w.warehouse_name,
		       o.doc_date, o.doc_type, o.total_amount, o.discount_amount, o.net_amount, o.vat_amount, o.grand_total,
		       o.update_by, o.update_date, o.is_active
		FROM tb_order o
		LEFT JOIN tb_customer c ON o.customer_id = c.customer_id
		LEFT JOIN tb_warehouse w ON o.warehouse_id = w.warehouse_id
		WHERE o.order_id = @ID AND o.is_delete = 0
	`
	var o models.Order
	var customerName, warehouseName *string
	err := db.QueryRow(queryHeader, sql.Named("ID", id)).Scan(
		&o.OrderID, &o.CustomerID, &customerName, &o.WarehouseID, &warehouseName,
		&o.DocDate, &o.DocType, &o.TotalAmount, &o.DiscountAmount, &o.NetAmount, &o.VatAmount, &o.GrandTotal,
		&o.UpdateBy, &o.UpdateDate, &o.IsActive,
	)
	if err == sql.ErrNoRows {
		return c.Status(404).SendString("Order not found")
	} else if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	o.CustomerName = customerName

	// 2. Items
	queryItems := `
		SELECT oi.auto_id, oi.order_id, oi.product_id, p.product_name,
		       oi.qty, oi.unit_price, oi.discount_amount, oi.line_total, oi.remark
		FROM tb_order_item oi
		LEFT JOIN tb_product p ON oi.product_id = p.product_id
		WHERE oi.order_id = @ID AND oi.is_delete = 0
	`
	rows, err := db.Query(queryItems, sql.Named("ID", id))
	if err != nil { return c.Status(500).SendString(err.Error()) }
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var productName *string
		if err := rows.Scan(&item.AutoID, &item.OrderID, &item.ProductID, &productName,
			&item.Qty, &item.UnitPrice, &item.DiscountAmount, &item.LineTotal, &item.Remark); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		item.ProductName = productName
		items = append(items, item)
	}

	return c.JSON(fiber.Map{"header": o, "items": items})
}

// InsertOrder เพิ่มใบสั่งขาย
func InsertOrder(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)
	
	type InsertRequest struct {
		Header models.Order      `json:"header"`
		Items  []models.OrderItem `json:"items"`
	}

	var req InsertRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString("Invalid Body")
	}

	tx, err := db.Begin()
	if err != nil { return c.Status(500).SendString(err.Error()) }

	user := "System" // Or from auth

	// 1. Header
	queryHeader := `
		INSERT INTO tb_order (
			customer_id, warehouse_id, doc_date, doc_type,
			total_amount, discount_amount, net_amount, vat_amount, grand_total,
			update_by, is_active
		) VALUES (
			@CustomerID, @WarehouseID, @DocDate, @DocType,
			@Total, @Discount, @Net, @Vat, @Grand,
			@UpdateBy, 1
		)
	`
	_, err = tx.Exec(queryHeader,
		sql.Named("CustomerID", req.Header.CustomerID),
		sql.Named("WarehouseID", req.Header.WarehouseID),
		sql.Named("DocDate", req.Header.DocDate),
		sql.Named("DocType", req.Header.DocType),
		sql.Named("Total", req.Header.TotalAmount),
		sql.Named("Discount", req.Header.DiscountAmount),
		sql.Named("Net", req.Header.NetAmount),
		sql.Named("Vat", req.Header.VatAmount),
		sql.Named("Grand", req.Header.GrandTotal),
		sql.Named("UpdateBy", user),
	)
	if err != nil { tx.Rollback(); return c.Status(500).SendString("Insert Header Fail: "+err.Error()) }

	// 2. Fetch ID
	var newID string
	err = tx.QueryRow("SELECT TOP 1 order_id FROM tb_order WHERE update_by = @UpdateBy ORDER BY autoID DESC", sql.Named("UpdateBy", user)).Scan(&newID)
	if err != nil { tx.Rollback(); return c.Status(500).SendString("Fetch ID Fail: "+err.Error()) }

	// 3. Items
	queryItem := `
		INSERT INTO tb_order_item (order_id, product_id, qty, unit_price, discount_amount, line_total, remark)
		VALUES (@OrderID, @ProductID, @Qty, @Price, @Discount, @Total, @Remark)
	`
	for _, item := range req.Items {
		_, err := tx.Exec(queryItem,
			sql.Named("OrderID", newID),
			sql.Named("ProductID", item.ProductID),
			sql.Named("Qty", item.Qty),
			sql.Named("Price", item.UnitPrice),
			sql.Named("Discount", item.DiscountAmount),
			sql.Named("Total", item.LineTotal),
			sql.Named("Remark", item.Remark),
		)
		if err != nil { tx.Rollback(); return c.Status(500).SendString("Insert Item Fail: "+err.Error()) }
	}

	if err := tx.Commit(); err != nil { return c.Status(500).SendString("Commit Fail: "+err.Error()) }

	return c.JSON(fiber.Map{"status": "success", "id": newID})
}

// UpdateOrderByID (Optional updates)
func UpdateOrderByID(c *fiber.Ctx) error {
	return c.Status(501).SendString("Update logic omitted for MVP")
}

// DeleteOrderByID (Soft)
func DeleteOrderByID(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)
	id := c.Params("id")
	
	tx, err := db.Begin()
	if err != nil { return c.Status(500).SendString(err.Error()) }

	// Header
	_, err = tx.Exec(`UPDATE tb_order SET is_delete = 1, is_active = 0 WHERE order_id = @ID`, sql.Named("ID", id))
	if err != nil { tx.Rollback(); return c.Status(500).SendString(err.Error()) }

	// Items
	_, err = tx.Exec(`UPDATE tb_order_item SET is_delete = 1 WHERE order_id = @ID`, sql.Named("ID", id))
	if err != nil { tx.Rollback(); return c.Status(500).SendString(err.Error()) }

	tx.Commit()
	return c.JSON(fiber.Map{"status": "success"})
}

// RemoveOrderByID (Hard)
func RemoveOrderByID(c *fiber.Ctx) error {
	db := c.Locals("db").(*sql.DB)
	id := c.Params("id")
	
	tx, err := db.Begin()
	if err != nil { return c.Status(500).SendString(err.Error()) }

	// Items first
	_, err = tx.Exec("DELETE FROM tb_order_item WHERE order_id = @ID", sql.Named("ID", id))
	if err != nil { tx.Rollback(); return c.Status(500).SendString(err.Error()) }

	// Header
	_, err = tx.Exec("DELETE FROM tb_order WHERE order_id = @ID", sql.Named("ID", id))
	if err != nil { tx.Rollback(); return c.Status(500).SendString(err.Error()) }

	tx.Commit()
	return c.JSON(fiber.Map{"status": "success"})
}
