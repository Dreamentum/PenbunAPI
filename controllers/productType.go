package controllers

import (
	"PenbunAPI/config"
	"PenbunAPI/models"
	"PenbunAPI/utils"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Select All (JOIN group_name)
func SelectAllProductType(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.type_group_name,               -- เก็บ group_code
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type p
		LEFT JOIN tb_product_type_group g
			ON g.group_code = p.type_group_name AND g.is_delete = 0
		WHERE p.is_delete = 0`)
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var out []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.TypeNameGroupCode, // ← ปรับชื่อฟิลด์
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		out = append(out, m)
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: out})
}

// Select Paging (JOIN group_name)
func SelectPageProductType(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.type_group_name,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type p
		LEFT JOIN tb_product_type_group g
			ON g.group_code = p.type_group_name AND g.is_delete = 0
		WHERE p.is_delete = 0
		ORDER BY p.update_date DESC
		OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`,
		sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var items []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.TypeNameGroupCode, // ← ปรับชื่อฟิลด์
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		items = append(items, m)
	}

	var total int
	if err := config.DB.QueryRow(`SELECT COUNT(*) FROM tb_product_type WHERE is_delete = 0`).Scan(&total); err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to count", Data: nil})
	}

	return c.JSON(models.ApiResponse{
		Status: "success", Message: "",
		Data: fiber.Map{"page": page, "limit": limit, "total": total, "items": items},
	})
}

// Select By ID (JOIN group_name)
func SelectProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	row := config.DB.QueryRow(`
		SELECT  p.product_type_id,
				p.type_name,
				p.type_group_name,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type p
		LEFT JOIN tb_product_type_group g
			ON g.group_code = p.type_group_name AND g.is_delete = 0
		WHERE p.product_type_id = @ID AND p.is_delete = 0`, sql.Named("ID", id))

	var m models.ProductType
	var upd sql.NullTime
	if err := row.Scan(
		&m.ProductTypeID,
		&m.TypeName,
		&m.TypeNameGroupCode, // ← ปรับชื่อฟิลด์
		&m.GroupName,
		&m.Description,
		&m.UpdateBy,
		&upd,
		&m.IDStatus,
		&m.IsDelete,
	); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "Not found", Data: nil})
		}
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
	}
	if upd.Valid {
		t := upd.Time
		m.UpdateDate = &t
	}

	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: m})
}

// Select By Name (LIKE)  (JOIN group_name)
func SelectProductTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	rows, err := config.DB.Query(`
		SELECT  p.product_type_id,
				p.type_name,
				p.type_group_name,
				ISNULL(g.group_name, '') AS group_name,
				p.description,
				p.update_by,
				p.update_date,
				p.id_status,
				p.is_delete
		FROM tb_product_type p
		LEFT JOIN tb_product_type_group g
			ON g.group_code = p.type_group_name AND g.is_delete = 0
		WHERE p.type_name LIKE '%' + @Name + '%'
		  AND p.is_delete = 0
		ORDER BY p.type_name ASC`, sql.Named("Name", name))
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to fetch", Data: nil})
	}
	defer rows.Close()

	var out []models.ProductType
	for rows.Next() {
		var m models.ProductType
		var upd sql.NullTime
		if err := rows.Scan(
			&m.ProductTypeID,
			&m.TypeName,
			&m.TypeNameGroupCode, // ← ปรับชื่อฟิลด์
			&m.GroupName,
			&m.Description,
			&m.UpdateBy,
			&upd,
			&m.IDStatus,
			&m.IsDelete,
		); err != nil {
			log.Println(err)
			return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to read", Data: nil})
		}
		if upd.Valid {
			t := upd.Time
			m.UpdateDate = &t
		}
		out = append(out, m)
	}
	if len(out) == 0 {
		return c.Status(404).JSON(models.ApiResponse{Status: "error", Message: "No matching product type found", Data: nil})
	}
	return c.JSON(models.ApiResponse{Status: "success", Message: "", Data: out})
}

// Insert — Thin (ให้ DEFAULT/Trigger จัดการเวลา/สถานะ)
func InsertProductType(c *fiber.Ctx) error {
	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				INSERT INTO tb_product_type (type_name, type_group_name, description, update_by)
				VALUES (@TypeName, @TypeNameGroupCode, @Description, @UpdateBy)`,
				sql.Named("TypeName", m.TypeName),
				sql.Named("TypeNameGroupCode", m.TypeNameGroupCode), // ← ใช้ชื่อตัวแปรใหม่
				sql.Named("Description", m.Description),
				sql.Named("UpdateBy", m.UpdateBy),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to insert",
			Data:    fiber.Map{"type_name": m.TypeName},
		})
	}
	return c.Status(201).JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type added successfully",
		Data:    fiber.Map{"type_name": m.TypeName},
	})
}

// Update — ให้ Trigger อัปเดตเวลาเอง
func UpdateProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var m models.ProductType
	if err := c.BodyParser(&m); err != nil {
		return c.Status(400).JSON(models.ApiResponse{Status: "error", Message: "Invalid request body"})
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				UPDATE tb_product_type SET
					type_name       = COALESCE(NULLIF(@TypeName,''),            type_name),
					type_group_name = COALESCE(NULLIF(@TypeNameGroupCode,''),   type_group_name),
					description     = @Description,
					update_by       = @UpdateBy
				WHERE product_type_id = @ID AND is_delete = 0`,
				sql.Named("TypeName", m.TypeName),
				sql.Named("TypeNameGroupCode", m.TypeNameGroupCode), // ← ชื่อใหม่
				sql.Named("Description", m.Description),
				sql.Named("UpdateBy", m.UpdateBy),
				sql.Named("ID", id),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{
			Status:  "error",
			Message: "Failed to update",
			Data:    fiber.Map{"product_type_id": id},
		})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type updated successfully",
		Data:    fiber.Map{"product_type_id": id},
	})
}

// Soft Delete — ให้ Trigger อัปเดตเวลาเอง
func DeleteProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Query("user")
	if username == "" {
		username = "UNKNOWN"
	}
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				UPDATE tb_product_type
				SET is_delete = 1, update_by = @UpdateBy
				WHERE product_type_id = @ID`,
				sql.Named("ID", id), sql.Named("UpdateBy", username),
			)
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to soft delete"})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type marked as deleted",
		Data:    fiber.Map{"product_type_id": id, "update_by": username},
	})
}

// Hard Delete
func RemoveProductTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	err := utils.ExecuteTransaction(config.DB, []func(tx *sql.Tx) error{
		func(tx *sql.Tx) error {
			_, err := tx.Exec(`DELETE FROM tb_product_type WHERE product_type_id = @ID`, sql.Named("ID", id))
			return err
		},
	})
	if err != nil {
		log.Println(err)
		return c.Status(500).JSON(models.ApiResponse{Status: "error", Message: "Failed to hard delete"})
	}
	return c.JSON(models.ApiResponse{
		Status:  "success",
		Message: "Product type removed successfully",
		Data:    fiber.Map{"product_type_id": id},
	})
}
